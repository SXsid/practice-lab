#!/usr/bin/env bash

arr=(10 25 3 47 8 19)

sum=0
for i in $(seq 1  ${#arr[@]});do
  if [[ $((i%2 )) -ne 0 ]];then
    sum=$((sum+${arr[$((i-1))]}))
  fi
done
echo $sum
