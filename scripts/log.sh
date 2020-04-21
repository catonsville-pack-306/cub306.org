#!/bin/bash

LOG_MAIN=~/logs/cub306.org/http/access_log
LOG_DEV=~/logs/develop.cub306.org/http/access_log
log="${LOG_MAIN}"
log2="${LOG_DEV}"

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
    alt_path="${2}"
    echo "Top 20 Visitor Counts:"
    cat "${log_path}" "${alt_path}" \
        | grep ' 200 ' \
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
    alt_path="${2}"
    echo "Top 20 Page Counts:"
    cat "${log_path}" "${alt_path}" \
        | grep ' 200 ' \
        | grep -v '\.css' \
        | grep -v '\.js' \
        | grep -v '\.png' \
        | gawk 'BEGIN { FPAT="([^ ]+)|(\"[^\"]+\")|(\\[[^\\]]+\\])" } { print $5 }' \
        | sed 's/ HTTP\/1\.1//' \
        | sort \
        | uniq -c \
        | sort --numeric-sort --reverse \
        | head -n 20
}

email_counts ()
{
    log_path="${1}"
    alt_path="${2}"
    echo "Top 20 Email Counts"
    cat "${log_path}" "${alt_path}" \
        | grep -o '/images/PackLogo_.*\?\w\sHTTP' "${log_path}" \
        | sed 's/HTTP//' \
        | sort \
        | uniq -c \
        | sort --numeric-sort --reverse \
        | head -n 20
}

work()
{
    echo EOL
}

while getopts "cehl:L:ps" opt; do
    case $opt in
        c) visitor_counts "${log}" "${log2}" ;;
        e) email_counts "${log}" "${log2}" ;;
        p) page_counts "${log}" "${log2}" ;;
        s) echo '----' ;;
        h)
            help
            exit 0
            ;;
        l) log="${OPTARG}" ;;
        L) log2="${OPTARG}" ;;
        *)
            printf "invalid options\n"
            help
            exit 1
            ;;
    esac
done

work
