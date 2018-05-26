#!/bin/bash

pushd /home/cubpack/cub306.org > /dev/null
if [ -n "$(git fetch --dry-run)" ] ; then
    git pull
fi
popd > /dev/null

