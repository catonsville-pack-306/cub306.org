#!/usr/bin/env ruby

#connect the CGI script to standard ruby paths
ENV['GEM_PATH'] = '/home/cubpack/.gems'
$LOAD_PATH.unshift("./public/cgi-bin")

#local
require "MarkdownToHTML"

mth = MarkdownToHTML.new
puts mth.render()
