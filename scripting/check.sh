#!/usr/bin/env bash
DirPath=$1

if [ -d "$DirPath" ]; then
    FileCount=$(find "$DirPath" -type f| wc -l)
    echo "size is $(du -sh "$DirPath"|cut -f1) containing $FileCount files"
    if [ "$FileCount" -gt 100 ] ;then
        echo "Large folder, be careful"
    fi
    
else
    echo " $DirPath path doesnt' exist"
    exit 1
fi


