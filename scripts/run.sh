#!/bin/bash

export DOCUMENT_ROOT=$(pwd)/public
export HTTP_ACCEPT="text/html"
export PATH_INFO="/index.md"
export REQUEST_URI="/index.md#section1"

export local_mode="on"

ruby public/cgi-bin/md.cgi
