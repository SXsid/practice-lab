
#!/usr/bin/env bash
# first if dosn't exist it has to equate to 0 for
# if marks is 300 then the count will ++ that too school id wise 
#

awk '
BEGIN{FS=","}
{
if (!($1 in freq)) freq[$1]=0;
if ($4>300) freq[$1]+=1;
}
END{ for(id in freq) print id":"freq[id]}
' school.csv
