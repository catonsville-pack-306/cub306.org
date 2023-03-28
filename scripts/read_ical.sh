#!/bin/zsh

# Documentation of the commands needed for each of the apps life cycles

# run the app in a loop to watch for changes to output
function run(){
    run1
}

function run1(){
    watch --color \
        "cat events.ics | \
        go run ReadIcal.go -date '2018-08-14'"
}

function run2(){
    watch --color \
        "cat events.ics | \
        grep -v 'X-APPLE' | \
        go run ReadIcal.go -date '2018-08-14'"
}


# build for linux hosts, it is assumed one can build for local use with no help
function build(){
    GOOS="$1" GOARCH="$2" go build -o readical.linux ReadIcal.go
}

# Run the tests in a loop and look for errors as code is written
function check(){
    watch --color go test ReadIcal.go ReadIcal_test.go
}

# scp the file up to a remote host
function send(){
    user=$1
    remote=$2
    scp ReadIcal $user@$remote:/home/$user/readical.linux
}

# ##############################################################################

function usage(){
    printf "Written by thomas.cherry@gmail.com\n"
    printf "A command to convert a time slice of ical events to markdown or HTML.\n"
    printf "\nUsage: -A <arch> -O <os> -R <name> -U <user> [-b -i -r -t]\n"
    printf "\tCommands executed in calling order.\n\n"
    template="%4s %-4s %15s %8s %-20s\n"

    printf "$template" "Flag" "Arg" "Default" "Name" "Description"
    printf "$template" "----" "----" "---------------" "------" "--------------"
    printf "$template" "-b" "" "" "build" "Create a linux binary"
    printf "$template" "-h" "" "" "help" "Display this text"
    printf "$template" "-i" "" "" "install" "Install binary to remote"
    printf "$template" "-r" "" "" "run" "Run app in a watch"
    printf "$template" "-t" "" "" "test" "Test app in a watch"
    printf "\n"
    printf "$template" "-A" "arch" "$goarch" "goarch" "Set GOARCH"
    printf "$template" "-O" "os" "$goos" "goos" "Set GOOS"
    printf "$template" "-R" "name" "$remote" "remote" "Set name for remote"
    printf "$template" "-U" "user" "$user" "user" "Set user for remote"
}

user=$USER
remote=hostname
goos=linux
goarch=amd64

while getopts bhirtA:O:U:R: opt ; do
    case "${opt}" in   
        b) build "$goos" "$goarch" ;;
        h) usage ; exit 0 ;;
        i) send "$user" "$remote" ;;
        r) run ;;
        t) check ;;

        A) goarch=${OPTARG};;
        O) goos=${OPTARG};;
        U) user=${OPTARG} ;;
        R) remote=${OPTARG} ;;

        *) usage ; exit 1 ;;
    esac
done
