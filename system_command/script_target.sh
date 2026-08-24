#!/usr/bin/env bash

read file
if [[ ! -e $file ]];then
  echo "NOT FOUND"
  exit 1
fi
if [[ ! -x $file ]];then
    echo "NOT EXECUTABLE"
    exit 2
else
    echo OK
    exit 0
  fi



