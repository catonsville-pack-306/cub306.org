#!/bin/bash

printf '<!DOCTYPE html>\n<html lang="en">\n<head>\n\t<meta charset="utf-8"/>\n'
printf '\t<title>%s</title>\n' "$(basename $(pwd))"
printf '</head>\n<body>\n\t<ol>\n'

list=$(find . -name '*.*' -print | grep -v '^.$' | grep -v 'index.html' | sort -r)

for i in $list
do
    printf "\t\t<li><a href="%s">%s</a></li>\n" $i $i
done

printf '\t</ol>\n</body>\n</html>'
