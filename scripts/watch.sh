#!/bin/bash

watch 'scripts/run.sh | tee public/out.html | lynx --stdin --dump'

