# cub306.org #
The public web site for pack 306 of Catonsville Maryland.

maintained by [Thomas Cherry](mailto:thomas.cherry@gmail.com)

## Overview ##

The site is constructed around an Apache [action](https://httpd.apache.org/docs/2.4/mod/mod_actions.html)
which executes a [CGI](https://help.dreamhost.com/hc/en-us/articles/217297307-CGI-overview)
script that coverts [Markdown](https://en.wikipedia.org/wiki/Markdown)
to HTML via the [Redcarpet gem](https://github.com/vmg/redcarpet).

## Writing Pages ##

Apache of course will deliver HTML files as normal, however, with this setup one
can also write Markdown files. All Markdown files will be converted to HTML and
embedded into the `default.erb` file. This file is an HTML file with Ruby
[ERB](https://ruby-doc.org/stdlib-2.5.1/libdoc/erb/rdoc/ERB.html) statements. The
basic web site layout will be handled by this file leaving the content to be handled
by the Markdown file. Each Markdown file can optionally have a matching ERB file
to handle more complicated layout issues. Both files must share the same name, only
changing by extension.

* default.erb
    * most_files.md
* fancy.erb
    * fancy.md

### Other Features ###

ERB files can include markdown files with:
    
    <%=inject "/alerts.md"%>

Explicit HTML title setting by adding the following to the top of the page:

    <!-- Title: Title_goes_here -->

If HTML title not set by title comment, then the first H1 header (single hash)
in the first 5 lines is the HTML title.
    
    # Title_goes_here #

Can add current ISO date and time with the "now" function

    <!= now %>

## Running on a Mac ##

Make the following changes to the apache config file :

    ...
    LoadModule cgi_module libexec/apache2/mod_cgi.so
    LoadModule actions_module libexec/apache2/mod_actions.so
    ...
    #ScriptAliasMatch...    #yes remove this
    ...
    AddHandler cgi-script .cgi
    ...
    Include /private/etc/apache2/vhosts/*.conf
    
Then create and edit the following file:
    sudo vim /private/etc/apache2/vhosts/cub306.site.conf

Then enter in the following, change as needed:
    ScriptAlias /cgi-bin/ "/Users/thomas/Documents/src/project/catonsville-pack-306/cub306.org/public/cgi-bin/"

    <Directory "/Users/thomas/Documents/src/project/catonsville-pack-306/cub306.org/public">
        AllowOverride None
        Options +ExecCGI
        AddHandler cgi-script .cgi
        Order allow,deny
        Allow from all
        Require all granted
    </Directory>

    <virtualHost *:80>
        DocumentRoot "/Users/thomas/Documents/src/project/catonsville-pack-306/cub306.org/public"
        ServerName cub306.local
        ErrorLog "/Users/thomas/Sites/var/log/cub306.site.error_log"

        <Directory "/Users/thomas/Documents/src/project/catonsville-pack-306/cub306.org/public">
            Options +Includes +FollowSymLinks +Indexes +ExecCGI
            AllowOverride All 
            Require all granted
            AddType type/html .shtml
            AddOutputFilter INCLUDES .shtml
        </Directory>
    </VirtualHost>
    
Finally restart apache with `sudo /usr/sbin/apachectl`


Apache may not find the gems installed for your user account, so you can try to update the version of ruby for the entire operating system (Macintosh in this case) as so:

    cd /System/Library/Frameworks/Ruby.framework/Versions/2.3/usr/bin
    sudo ./gem install redcarpet
    sudo ./gem install rqrcode
    
In other cases, such as on a hosting service, you may nat be able to install gems in this way. For these cases, you will need to use a "wrapper" such as [qrc.cgi](public/cgi-bin/qrc.cgi). Here the idea is to set the gem path and then call ruby:

	#!/bin/bash
	GEM_HOME=/home/<account_name_here>/.gems
	export GEM_HOME
	./qrcode.cgi

After making changes on a hosted system, you may need to run `touch ~/cub306.org/public/tmp/restart.txt` to get "[Passenger](https://help.dreamhost.com/hc/en-us/articles/215769578-Passenger-overview)" to recognize the change.
