
#!/usr/bin/env bash

awk 'BEGIN{FS=","} { if (max[$1]<$4) {max[$1]=$4 ; schools[$1]=$2 }} END{ for (school in schools) print schools[school]}'  marks.csv
