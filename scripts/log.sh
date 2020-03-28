#!/bin/bash

log=~/logs/cub306.org/http

help()
{
    printf "Update repository on demand\n"
    printf "cmd [dmh]\n"
    t1="%5s %s\n"
    printf "${t1}" "-d" "development"
    printf "${t1}" "-h" "help"
    printf "${t1}" "-m" "master"
}

visitor_counts ()
{
    log_path=$1
    grep ' 200 ' "${log_path}" \
        | awk '{print $1}' \
        | sort \
        | uniq -c
}

page_counts ()
{
    log_path=$1
    #grep ' 200 ' "${log_path}"
    grep ' 200 ' "${log_path}" \
        | awk 'BEGIN { FPAT="([^ ]+)|(\"[^\"]+\")|(\\[[^\\]]+\\])" } { print $5 }' \
        | sort \
        | uniq -c
}


work()
{
    echo done
}

while getopts "chl:ps" opt; do
    # OPTARG - not used yet
    case $opt in
        s) echo '----' ;;
        p) page_counts "${log}" ;;
        c) visitor_counts "${log}" ;;
        h)
            help
            exit 0
            ;;
        l) log="${OPTARG}" ;;
        *)
            printf "invalid options\n"
            help
            exit 1
            ;;
    esac
done

work
