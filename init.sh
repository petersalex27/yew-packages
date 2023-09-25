#!/bin sh

if [ -z $1 ]; then
    echo "usage: sh init.sh <dir> [name]"
    exit
fi

modpath="github.com/petersalex27/yew-packages"
d=$(pwd)

# get module name
mod=""
if [ -z $2 ]; then
    if [ "$1" = "." ]; then
        mod=""
    else
        mod="/$1"
    fi
else
    if [ "$2" = "." ]; then
        mod=""
    else
        mod="/$2"
    fi
fi

if cd ./$1; then
    if [ -f "go.mod" ]; then
        rm -f go.mod
    fi

    go mod init $modpath$mod
    cd $d
else
    echo no directory named $1
fi