# built in
require 'date'      # iso date method
require 'erb'       # template pages and inject pages can be ERB files

# external
require 'redcarpet' # gem install redcarpet

# CGI to process markdown into HTML
class MarkdownToHTML
    # object initialization, will read the following envirnment variables:
    # * DOCUMENT_ROOT - path on server to where the public directory is
    # * PATH_INFO - file requested on the url starting with / and before ?
    # * REQUEST_URI - file request starting with / and including ? & # sections
    # * HTTP_ACCEPT - standard HTTP Accept header
    # * local_mode - if "on" render output will be HTML only, no headers
    def initialize
        # http://www.cgi101.com/book/ch3/text.html
        # https://www.nginx.com/resources/wiki/start/topics/examples/phpfcgi/
        # http://lemp.test/test.php/foo/bar.php?v=1
        # 'DOCUMENT_ROOT' => '/var/www'
        # 'PATH_INFO' => '/test.php/foo/bar.php',
        # 'REQUEST_URI' => '/test.php/foo/bar.php?v=1'

        @root = read_env('DOCUMENT_ROOT', nil, %r{[^\w\/\.-]})
        @doc_uri = read_env('PATH_INFO', '/index.md', %r{[^-\/\w\?&=\.#%]*})
        @req_uri = read_env('REQUEST_URI', '/index.md?', %r{[^-\/\w\?&=\.#%]*})
        @accept = read_env('HTTP_ACCEPT', '*/*', %r{[^\/\w\+\.,;=: \*]})

        @file_base_name = File.basename(@req_uri)
        @file_path = File.path(@req_uri)
        @file_extname = File.extname(@req_uri)
        @file_dirname = File.dirname(@req_uri)

        # template = 'Content-type: <%=ctype%>; charset=utf-8\n\n<%=page%>'
        template = "Content-type: text/html; charset=utf-8\n\n<%=page%>"
        @template = ERB.new template
        @debug = false
        local_mode
    end

    # enable local mode which is used when testing on the CLI and not CGI
    def local_mode
        @template = ERB.new '<%=page%>' if ENV['local_mode'] == 'on'
    end

    # Read an envirnment variable and filter it out for safty
    # @param name envirnment variable to read
    # @param default value to use if there is no variable
    # @param filter regexp to use to filter out of the variable
    # @return filtered variable or the default value
    def read_env(name, default = nil, filter = /[^\w]*/)
        raw = ENV[name].nil? ? default : ENV[name]

        raw = raw.downcase.gsub(filter, '\1') unless raw.nil?
        raw
    end

    # public method, this action is what does the work of creating the CGI
    # output
    # @return CGI output, headers, then HTML
    def render
        # read in markdown
        contents = load_file(@root, @doc_uri)
        title = find_title(contents, File.basename(@doc_uri, '.*'))

        wrapper_name = find_wrapper_page
        page = markdown(contents)
        unless wrapper_name.nil?
            page = ERB.new(load_file(@root, wrapper_name)).result(binding)
        end

        # final CGI rendering
        @template.result(binding)
    end

    # extract a title from the content of markdown
    # @param contents markdown text
    # @param title default title to use, assumes nil
    # @return document title or the default value
    def find_title(contents, title = nil)
        requested_title = contents[/^\s*<!-- Title: (.*) -->\s*$/, 1]
        if !requested_title.nil?
            title = requested_title
        else
            requested_title = contents.lines.to_a[0, 5].join[/^# (.*) #\s*$/, 1]
            title = requested_title unless requested_title.nil?
        end
        title
    end

    # find a matching wrapper page on the file sytem
    # @param root where to start looking, not ending in a slash
    # @param doc_url request URI not starting with a slash
    # @return URI of the matching erb file
    def find_wrapper_page(root = @root, doc_uri = @doc_uri)
        wrapper_name = nil
        # find custome wrapper page
        if File.exist?("#{root}/#{doc_uri}.erb")
            wrapper_name = "/#{doc_uri}.erb"
        elsif File.exist?("#{root}/default.erb")
            wrapper_name = '/default.erb'
        end
        wrapper_name
    end

    # Converts markdown to html
    # @param markdown markdown text
    # @return string - html
    def markdown(markdown = '')
        render = Redcarpet::Render::HTML.new
        options = { tables: true, strikethrough: true, underline: true }
        Redcarpet::Markdown.new(render, options).render(markdown.dup)
    end

    # loads a file and returns it's content
    # @param path directory to find file
    # @param name file name
    # @return contents of file at path/name
    def load_file(path, name)
        handle = File.open("#{path}#{name}", 'r')
        text = handle.read
        text
    end

    # return the curent date and time
    # @return current time in iso like format
    def now
        Time.now.strftime('%Y-%m-%d %H:%M:%S')
    end

    # loads a markdown file, translates it, processes ERB commands
    # @param path URI for the markdown file to load
    # @return HTML content
    def inject(path)
        return "#{path} does not exist" unless File.exists?(@root + path)
        
        contents = load_file(@root, path)
        
        # try first
        wrapper_erb = ERB.new contents
        contents = wrapper_erb.result(binding)
        
        page = markdown(contents)

        wrapper_erb = ERB.new page
        result = wrapper_erb.result(binding)

        result
    end
    
    def blog_exists
        if @req_uri.end_with? '.md'
            if @req_uri.end_with? 'index.md'
                return true
            end
            return false
        end
        uri = @req_uri
        File.directory?(@root + uri + "blog")
    end
    
    def blog_here
        blog(@file_dirname + "/blog")
    end
    
    def blog_path file
        @file_dirname + "/blog/" + file
    end
    
    def blog(path)
        result = ""
        todays_file = Time.now.strftime("%Y-%m-%d")

        Dir[@root + path + "/*.md"].each do |full_path|
            file_name = File.basename(full_path)
            file_name_no_ext = File.basename(full_path, ".*")
            parts = file_name_no_ext.split("-to-")
            inject_file = path + "/" + file_name
            
            begin
                if parts.count == 1 && todays_file==file_name_no_ext
                    result = result + helper_wrap_in_article( inject_file )
                elsif parts.count == 2
                    start = Date.parse(parts[0])
                    stop = Date.parse(parts[1])
                    if Date.today.between?(start, stop)
                        result = result + helper_wrap_in_article( inject_file )
                    end
                end
            rescue StandardError => e
                result = e
            end
        end
        result
    end
    
    def helper_wrap_in_article(inject_file)
        c = inject(inject_file)
        inject_file["//"] = "/"
        %Q(
<article>
    #{c}
    <div class="tools">
        <a href="#{inject_file}"><i class="fas fa-link"></i></a>
    </div>
</article>) unless c.length < 1
    end

end
