#!/usr/bin/env ruby

ENV['GEM_PATH'] = '/home/cubpack/.gems'

#built in
require 'date'
#require 'file'
require "erb"

#external
require 'redcarpet'	#gem install redcarpet

class MarkdownToHTML
    def initialize()
        <<-DOC
        http://www.cgi101.com/book/ch3/text.html
        https://www.nginx.com/resources/wiki/start/topics/examples/phpfcgi/
        
        http://lemp.test/test.php/foo/bar.php?v=1
        
        /home/cubpack/cub306.org/public
        'DOCUMENT_ROOT' => '/var/www'
        'DOCUMENT_URI' => '/test.php/foo/bar.php',
        'REQUEST_URI' => '/test.php/foo/bar.php?v=1'
        DOC
        
        @root = read_env("DOCUMENT_ROOT")
        @doc_uri = read_env("DOCUMENT_URI", "/index.md")
        @req_uri = read_env("REQUEST_URI", "/index.md")
        @accept = read_env("HTTP_ACCEPT", "*/*")
        
        @template=ERB.new "Content-type: <%=ctype%>; charset=utf-8\n\n<%=page%>"
        @debug = false
        local_mode()
    end
    
    def local_mode()
        if ENV["local_mode"]=="on"
            @template=ERB.new "<%=page%>"
        end
    end
    
    def read_env(name, default=nil)
        if ENV[name].nil? ; default else ENV[name] end
    end
    
    def render()
        ctype = "text/html" #@accept
        page = ""
        title = File.basename(@doc_uri, ".*")
        now = DateTime.now.strftime("%Y-%m-%d %H:%M:%S")
        
        #read in markdown
        contents = load_file(@root, @doc_uri)
        page = markdown(contents)
           
        #read in default wrapper
        wrapper_name = nil
        
        #find wrapper page
        untyped_title = @doc_uri.gsub(/\.md$/, "\.erb")
        if File.exist?("#{@root}/#{untyped_title}")
            wrapper_name = "/#{untyped_title}"
        elsif File.exist?("#{@root}/default.erb")
            wrapper_name = "/default.erb"
        end
        unless wrapper_name.nil?
            wrapper_text = load_file(@root, wrapper_name)
            wrapper_erb = ERB.new wrapper_text
            page = wrapper_erb.result(binding)
        end
        
        #final CGI rendering
        ret = @template.result(binding)
        ret
    end
    
    def markdown(markdown = "")
        render = Redcarpet::Render::HTML.new
        options = {tables:true, strikethrough:true, underline:true}
        Redcarpet::Markdown.new(render, options).render(markdown.dup)
    end
    
    def load_file(path, name)
        handle = File.open(path + name, "r")
        text = handle.read
        text
    end
    #    <%= ERB.new(File.read(File.expand_path("#{@root}/alerts.md")), nil, nil, '_sub01').result(binding) %>

    def sub_page(path)
        return ERB.new(File.read(File.expand_path(path)), nil, nil, '_sub01').result(binding)
    end
end

mth = MarkdownToHTML.new
puts mth.render()
