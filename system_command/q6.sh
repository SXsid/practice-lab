
#!/usr/bin/env bash
#INFO: cocnept NF give the number of words after seprator applied
awk '
{
if (NF==2) print $2,$1;
else print $0;
}' swap.txt
