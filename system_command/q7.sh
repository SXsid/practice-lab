#!/usr/bin/env bash

read n
if [[ $n =~ ^[0-9]+$ && $n -ne "0" ]];then
  mkdir output_files
  for i in $(seq 1 "$n");do
    touch output_files/file_${i}.txt
  done
else
  echo "Error $n is not a positive number"
  exit 1
fi
