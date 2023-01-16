#!/bin/zsh

function test(){
	watch --color go test ReadIcal.go ReadIcal_test.go
}

function run(){
	watch --color cat events.ics | \
        grep -v 'X-APPLE' | \
		go run ReadIcal.go -date '2018-08-14'
}

while getopts etr opt ; do
    case "${opt}" in   
        t) test ;;
        r) run ;;
        e) echo ${OPTARG} ;;
        
        *) usage ;;
    esac
done
