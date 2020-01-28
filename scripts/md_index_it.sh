#!/bin/bash

printf '## %s\n\n' "$(basename $(pwd))"

list=$(find . -name '*.*' -print | grep -v '^.$' | grep -v 'index.html' | sort -r)

for i in $list
do
    printf "* [%s](%s)\n" $i $i
done

printf '\n'
