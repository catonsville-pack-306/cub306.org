#!/bin/bash

if [ -n "$(git fetch --dry-run)" ] ; then
    cd /home/cubpack/cub306.org
    git pull
fi

