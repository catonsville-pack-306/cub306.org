#!/bin/bash

export DOCUMENT_ROOT=$(pwd)/public
export HTTP_ACCEPT="text/html"

export local_mode="on"

build()
{
    export PATH_INFO="/${2}.md"
    export REQUEST_URI="/${2}.md#section1"

    ruby public/cgi-bin/md.cgi | sed 's/.md/\.html/'> ${1}/${2}.html
}

if [ ! -d out ] ; then mkdir out ; fi

for i in $(ls public/*.md) ; do
    name=$(basename -s .md $i)
    build out $name
done

cp -r public/images out/.
cp -r public/scripts out/.
cp -r public/styles out/.

tree out