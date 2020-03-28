#!/bin/bash

LOG_MAIN=~/logs/cub306.org/http/access_log
LOG_DEV=~/logs/develop.cub306.org/http/access_log
log="${LOG_MAIN}"

help()
{
    printf "Dump some simple log analyzes\n"
    printf "cmd [dmh]\n"
    t1="%5s %s\n"
    printf "${t1}" "-d" "development"
    printf "${t1}" "-h" "help"
    printf "${t1}" "-m" "master"
}

visitor_counts ()
{
    log_path="${1}"
    grep ' 200 ' "${log_path}" \
        | grep -v '\.css' \
        | grep -v '\.js' \
        | awk '{print $1}' \
        | sort \
        | uniq -c \
        | sort --numeric-sort --reverse \
        | head -n 20
}

page_counts ()
{
    log_path="${1}"
    grep ' 200 ' "${log_path}" \
        | grep -v '\.css' \
        | grep -v '\.js' \
        | gawk 'BEGIN { FPAT="([^ ]+)|(\"[^\"]+\")|(\\[[^\\]]+\\])" } { print $5 }' \
        | sort \
        | uniq -c \
        | sort --numeric-sort --reverse \
        | head -n 20
}

email_counts ()
{
    log_path="${1}"
    grep -o '/images/PackLogo_.*\?\w\sHTTP' "${log_path}" \
        | sed 's/HTTP//' \
        | sort \
        | uniq -c \
        | sort --numeric-sort --reverse \
        | head -n 20

#grep 'GET /images/PackLogo_' ~/logs/cub306.org/https/access.log | sort | grep -o '/images/PackLogo_.*\?\w\sHTTP' | sed 's/HTTP//' |  uniq -c | sort --numeric-sort --reverse


}

work()
{
    echo done
}

while getopts "cehl:ps" opt; do
    # OPTARG - not used yet
    case $opt in
        c) visitor_counts "${log}" ;;
        e) email_counts "${log}" ;;
        p) page_counts "${log}" ;;
        s) echo '----' ;;
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
