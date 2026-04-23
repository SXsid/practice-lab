#! /usr/bin/env bash

greet() {
    Name=$1
    echo "hello $Name"

}

isFile() {
    FilePath=$2
    if [ -f "$FilePath" ];then
        echo "is a file"
    else
        echo "not a file"
        return 1
    fi
}
greet "sid"
greet "$@"
isFile "d" "./file.sh"


