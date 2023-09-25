#!/bin sh

if [ -z $1 ]; then
    echo "usage: sh clean.sh <dir>"
    exit
fi

d=$(pwd)

if cd ./$1; then
    if [ -f "go.mod" ]; then
        go clean -cache
        go clean -modcache
    else
        echo "directory does not contain a go module"
    fi
    cd $d
else
    echo no directory named $1
fi