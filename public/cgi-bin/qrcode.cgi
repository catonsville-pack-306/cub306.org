#!/usr/bin/env ruby

require 'cgi'
require 'rqrcode'

begin

protocol = ENV["Server_Protocol"].nil? ? "http" : ENV["Server_Protocol"]
server = ENV["SERVER_NAME"] == "" ? "" : ENV["SERVER_NAME"]
port = ENV["SERVER_PORT"] == "" ? "" : ENV["SERVER_PORT"]

port = "" if port = ":80"
port = ":"+port unless port == ""

base = protocol + "://" + server + port

name = "REQUEST_URI"
raw = ENV[name].nil? ? default : ENV[name]
raw = raw.downcase.gsub(%r{[^-/\/\w\?&=\.#%]*}, '\1') unless raw.nil?
#raw = raw.gsub("--", "/")
raw_params = raw.gsub(/^.*\?/, "")

params = CGI::parse raw_params

output_types = {
    "png" => "image/png",
    "html" => "text/html",
    "svg" => "image/svg+xml"
}

format_name = "png"
content_type = output_types[format_name]
out = ""
options = {size: 120}

note = nil

params.each do |key, value|
    case key
    when "url"
        out = value[0]
    when "path"
        out = base + "/" + value[0]
    when "format"
        format_name = value[0]
        raw = output_types[format_name]
        content_type = raw unless raw.nil?
    when "size"
        options['size'] = value[0].to_i
        #note = options["size"]
    end
end

qrcode = RQRCode::QRCode.new(out)

case content_type
when "image/png"
    #image = qrcode.as_png ( options )
    #image = qrcode.as_png ( {size: options['size'] } )
    image = qrcode.as_png size: options['size']
when "text/html"
    image = qrcode.as_html
when "image/svg+xml"
    image = qrcode.as_svg
when "text/ansi"
    image = qrcode.as_ansi
when "text/text"
    image = qrcode.to_s
end

if note.nil?
    puts "Content-type: #{format_name}\n\n#{image}"
else
    puts "Content-type: text/text\n\n#{note}"
end

rescue StandardError => e
    puts "Content-type: text/text\n\n#{e}"
end

################################################################################

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
