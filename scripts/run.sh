#!/bin/bash

export DOCUMENT_ROOT=$(pwd)/public
export HTTP_ACCEPT="text/html"
export DOCUMENT_URI="/index.md"
export REQUEST_URI="/index.md#section1"

ruby public/cgi-bin/md.rb
