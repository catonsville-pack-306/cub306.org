#!/usr/bin/env ruby

#connect the CGI script to standard ruby paths
ENV['GEM_PATH'] = '/home/cubpack/.gems'

if File.exist?("#{ENV['DOCUMENT_ROOT']}/cgi-bin")
    $LOAD_PATH.unshift("#{ENV['DOCUMENT_ROOT']}/cgi-bin")
end

#local
begin
    require "PathToQRCode"
rescue Exception => e
    puts "Content-type: text/text\n\n#{e}"
end

qrcode = PathToQRCode.new
puts qrcode.render()

#puts "Content-type: text/text\n\n#{ENV['DOCUMENT_ROOT']}/cgi-bin"
