#!/bin/bash
# searches through the battle maps stored on Cloudfront.
# $1 is the end date and is required.
# $2 is the start date and is optional. If not provided, it starts 2 days before the end date.
# start and end date formats are YYYYMMDD
# this tool stops at the first map it finds, searching from the end date towards the start date.
# On Linux/OSX this shouhld run from the command line fine.
# On Windows, you will need either the "Linux Sub System" installed
# with Bash support, or the Git-Bash tool.
END=${1}
ENDD=$(date -d@$END || date -r$END)

START=${2:-$(date -d"$ENDD -2 days" +%s || date -j -r$END -v-2d +%s)}

for ((i=$END;i>$START;i-=1)) ; do
#for i in $(seq -f %02g $END -1 $START); do
    URL="https://d2n1d3zrlbtx8o.cloudfront.net/download/BattleMap.zip/AreaMap_002.${i}"
    echo $URL
 #   exit
    CODE=$(curl -o /dev/null --silent --head --write-out '%{http_code}\n' "${URL}")
    echo "$i: $CODE"
    if [[ "200" -eq "$CODE" ]] ; then
        curl -o "AreaMap_002.${i}" "https://d2n1d3zrlbtx8o.cloudfront.net/download/BattleMap.zip/AreaMap_002.${i}"
        break
    fi
done
