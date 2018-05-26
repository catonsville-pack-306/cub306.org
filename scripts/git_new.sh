#!/bin/bash

pushd /home/cubpack/cub306.org
if [ -n "$(git fetch --dry-run)" ] ; then
    git pull
fi
popd

