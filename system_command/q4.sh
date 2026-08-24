
#!/usr/bin/env bash

awk '
BEGIN{FS=","}
{
  if (NR==1) min=$3;
if ($3<min){
min=$3
}
if (!($2 in dept)|| dept[$2]>$3){
dept[$2]=$3;
}
}
END{
print min;
for (d in dept){
print  d":"dept[d];
}
}
' emp.csv
