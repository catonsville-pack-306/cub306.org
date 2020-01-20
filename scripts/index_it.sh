#!/bin/bash

printf '<!DOCTYPE html>\n<html lang="en">\n<head>\n\t<meta charset="utf-8"/>\n'
printf '\t<title>%s</title>\n' "$(basename $(pwd))"
printf '</head>\n<body>\n\t<ol>\n'

find . -name '*.*' -exec printf "\t\t<li><a href="%s">%s</a></li>\n" {} {} \;

printf '\t</ol>\n</body>\n</html>'
