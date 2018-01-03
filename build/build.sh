#!/bin/bash -e
if [ "$#" == "3" ]; then
    REGISTRY=$1
    NAMESPACE=$2
    TAG=$3
fi

rm -f root/bin/healthcheckapp
GOOS=linux GOARCH=amd64 go build -o root/bin/healthcheckapp ../main.go
(cd root && fakeroot tar cvf ../root.tar .)
image=${REGISTRY}/${NAMESPACE}/healthcheckapp

docker build -t ${image}:${TAG} .

#docker push ${image}:${TAG}
#docker push ${image}:latest
rm root.tar
