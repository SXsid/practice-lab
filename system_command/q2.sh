#!/usr/bin/env bash

awk   'BEGIN{FS=","} {if ($4>map[$1]) map[$1]=$4 } END{for( school in map) print school,map[school]}' marks.csv
