#!/bin/bash

# Ruby does not let me set this, do it in bash first
GEM_HOME=/home/cubpack/.gems
export GEM_HOME

# Does not seam to be needed
#GEM_PATH=/home/cubpack/.gems
#export GEM_PATH

./qr.ruby.cgi
