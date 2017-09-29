#!/bin/bash

export GOPATH=$GOPATH
export GOBIN=$GOPATH/bin
CUR_DIR=$(pwd)
name=${CUR_DIR##*/}

echo "Building $name:"
go build -o $name $2

mkdir -p ./lib
# Make map of what libs actually exist in the busy box image.
declare -A libMap
baseid=$(docker create busybox:ubuntu-14.04 /bin/sleep 10000)
docker start $baseid
lsR=$(docker exec $baseid ls /lib)
while read -r line; do
    libMap[$line]=1
done <<< "$lsR"
docker rm -vf $baseid
# Get a list of the libs that the go binary references, and
# add the ones that aren't in the busy box libMap
ldd $name | while read -r line; do
    name=$(echo "$line" | awk -F "=>" '{print $1}' | sed -e 's/^[[:space:]]*//' -e 's/[[:space:]]*$//')
    pt=$(echo "$line" | awk -F "=>" '{print $2}')
    path=$(echo "$pt" | awk -F " " '{print $1}')
    if [[ $path =~ \(.* || $path = "" || ${libMap[$name]} = 1 ]]
    then
	continue
    fi
    cp $path ./lib/
done

echo "Building docker image:"
docker build -t daticahealth/$name:release .

echo "Cleaning up:"
rm -rf ./lib
