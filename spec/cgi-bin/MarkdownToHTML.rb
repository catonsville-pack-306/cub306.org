#require 'md.cgi'

#$LOAD_PATH.unshift("./public/cgi-bin")

require 'MarkdownToHTML'

describe MarkdownToHTML, "basic" do
    context "as an object" do
        it "has properties" do
            markdown = MarkdownToHTML.new
            #expect(markdown.instance_variable_get :@root).to eq nil
            expect(markdown.instance_variable_get :@doc_uri).to eq "/index.md"
            expect(markdown.instance_variable_get :@req_uri).to eq "/index.md?"
            expect(markdown.instance_variable_get :@accept).to eq "*/*"
        end
    end
    context "has a base path" do
        it "does not have a valid value" do
            markdown = MarkdownToHTML.new
            
            raw = markdown.read_env("test")
            expect(raw).to eq nil
            
            raw = markdown.read_env("test", "blank")
            expect(raw).to eq "blank"
        end
        it "is a valid pase path" do
            markdown = MarkdownToHTML.new
            
            ENV["test"] = "/some/path"
            raw = markdown.read_env ("test")
            puts raw
            expect(raw).to eq "somepath"
            
            raw = markdown.read_env("test", nil, filter=/[^\w\/]*/)
            expect(raw).to eq "/some/path"

            ENV["test"] = "/hello world - test_this_1?name=value&a=%20b#item"
            raw = markdown.read_env("test", nil, /[^\w\/\.-]/)
            expect(raw).to eq "/helloworld-test_this_1namevaluea20bitem"
            
            raw = markdown.read_env("test", nil, /[^\/\w\?&=\.#%]*/) #no '-' ?
            expect(raw).to eq "/helloworldtest_this_1?name=value&a=%20b#item"
            
            ENV["test"] = "Accept: text/html, application/xhtml+xml, application/xml;q=0.9, */*;q=0.8"
            raw = markdown.read_env("test", nil, /[^\/\w\+\.,;=: \*]/)
            expect(raw).to eq ENV["test"].downcase
        end
        it "now timer" do
            markdown = MarkdownToHTML.new
            raw = markdown.now
            expect(raw).to eq DateTime.now.strftime("%Y-%m-%d %H:%M:%S")
        end
        it "markdown test" do
            markdown = MarkdownToHTML.new
            raw = markdown.markdown("* bullet")
            expect(raw).to eq ("<ul>\n<li>bullet</li>\n</ul>\n")
            
            mtable = <<-DOC
| a | b |
| - | - |
| 1 | 2 |
| 3 | 4 |
DOC
            expectedTable = <<-DOC
<table><thead>\n<tr>\n<th>a</th>\n<th>b</th>\n</tr>\n</thead><tbody>
<tr>\n<td>1</td>\n<td>2</td>\n</tr>
<tr>\n<td>3</td>\n<td>4</td>\n</tr>
</tbody></table>
DOC
            raw = markdown.markdown(mtable)
            expect(raw).to eq (expectedTable)

            
            raw = markdown.markdown("~~strike~~")
            expect(raw).to eq ("<p><del>strike</del></p>\n")

            raw = markdown.markdown("_underline_")
            expect(raw).to eq ("<p><u>underline</u></p>\n")
            
        end
    end
end
