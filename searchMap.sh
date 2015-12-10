#!/bin/bash

END=${1}
ENDD=$(date -d@$END)

START=${2:-$(date -d"$ENDD -2 days" +%s)}

for i in $(seq $END -1 $START); do
    URL="https://d2n1d3zrlbtx8o.cloudfront.net/download/BattleMap.zip/AreaMap_002.${i}"
    CODE=$(curl -o /dev/null --silent --head --write-out '%{http_code}\n' "${URL}")
    echo "$i: $CODE"
    if [[ "200" -eq "$CODE" ]] ; then
        curl -o "AreaMap_002.${i}" "https://d2n1d3zrlbtx8o.cloudfront.net/download/BattleMap.zip/AreaMap_002.${i}"
        break
    fi
done
