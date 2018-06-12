#!/bin/bash

BRANCH=master
SITE=cub306.org
URL=https://github.com/catonsville-pack-306/cub306.org

help()
{
    printf "Update repository on demand\n"
    printf "cmd [dmh]\n"
    t1="%5s %s\n"
    printf "${t1}" "-d" "development"
    printf "${t1}" "-h" "help"
    printf "${t1}" "-m" "master"
}

work()
{
    if [ -a /home/cubpack/${SITE} ] ; then
        pushd /home/cubpack/${SITE} > /dev/null
        revision=$(git ls-remote $URL | grep ${BRANCH} | cut -f 1)
        if [ -n "${revision}" ] ; then
            # we got a revision number from curl
            if [ -f ~/.last_revision_of_${BRANCH} ] ; then
                old_revision=$(cat ~/.last_revision_of_${BRANCH})
                if [ "${old_revision}" != "${revision}" ] ; then
                    git pull
                fi
            else
                #first time here
                git pull
            fi
            echo ${revision} > ~/.last_revision_of_${BRANCH}
        fi
        popd > /dev/null
    fi
}

while getopts "dhm" opt; do
    # OPTARG - not used yet
    case $opt in
        d)
            BRANCH="develop"
            SITE="develop.cub306.org"
            ;;
        m)
            BRANCH="master"
            SITE="cub306.org"
            ;;
        h)
            help
            exit 0
            ;;
        *)
            printf "invalid options\n"
            help
            exit 1
            ;;
    esac
done

work
