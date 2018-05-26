#!/usr/bin/env ruby

ENV['GEM_PATH'] = '/home/cubpack/.gems'
require File.join(File.dirname(__FILE__), 'boot')

puts "Content-type: text/html; charset=utf-8\n\n"
puts "home=#{ENV['GEM_HOME']}<br>\n"
puts "path=#{ENV['GEM_PATH']}<br>\n"

