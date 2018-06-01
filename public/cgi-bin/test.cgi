#!/usr/bin/env ruby

ENV['GEM_PATH'] = '/home/cubpack/.gems'
#require File.join(File.dirname(__FILE__), 'boot')

puts "Content-type: text/html; charset=utf-8\n\n"
puts "ghome=#{ENV['GEM_HOME']}<br>\n"
puts "gpath=#{ENV['GEM_PATH']}<br>\n"

puts "home=#{ENV['HOME']}<br>\n"
puts "script=#{ENV['SCRIPT_NAME']}<br>\n"
puts "fpath=#{ENV['PATH_TRANSLATED']}<br>\n"

ENV.each do |n,v|
    $stderr.puts "#{n}=#{v}\n"
end
