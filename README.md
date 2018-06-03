# web_site
The public web site for pack 306

maintained by thomas.cherry@gmail.com


## Running on a Mac ##

Make the following changes to apache:

    LoadModule cgi_module libexec/apache2/mod_cgi.so
    #ScriptAliasMatch...
    AddHandler cgi-script .cgi
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
    
Finally restart apache
    sudo /usr/sbin/apachectl


Apache may not find the gems installed for your user account, so you can try to update the version of ruby for the entire operating system as so:

    cd /System/Library/Frameworks/Ruby.framework/Versions/2.3/usr/bin
    sudo ./gem install redcarpet