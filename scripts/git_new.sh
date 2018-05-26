#!/bin/bash

if [ "$(git status -uno | grep -o 'Your branch is up-to-date')" -ne "Your branch is up-to-date" ] ; then
    cd /home/cubpack/cub306.org
    git pull
fi

