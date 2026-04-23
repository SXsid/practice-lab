#! /usr/bin/env bash

FilePath=$1
if [ -z "$FilePath" ] ;then 
    echo "usage: ./file.sh <file name or folder path>"
    exit 1
    
fi

if [ -f "$FilePath" ]; then
    #execute a command and use it as output
    echo "Size is :$(du -sh "$FilePath" | cut -f1)"
    exit 0
else
    echo "$FilePath is not a file"
    exit 1
fi
