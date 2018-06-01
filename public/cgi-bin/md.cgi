#!/usr/bin/env ruby

#connect the CGI script to standard ruby paths
ENV['GEM_PATH'] = '/home/cubpack/.gems'

if File.exist?("#{ENV['DOCUMENT_ROOT']}/cgi-bin")
    $LOAD_PATH.unshift("#{ENV['DOCUMENT_ROOT']}/cgi-bin")
end

#local
require "MarkdownToHTML"

mth = MarkdownToHTML.new
puts mth.render()
