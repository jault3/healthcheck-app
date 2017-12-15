#!/bin/bash -e
rm -f root/bin/healthcheckapp
go build -o root/bin/healthcheckapp ../main.go
(cd root && fakeroot tar cvf ../root.tar .)
image=${REGISTRY}/${NAMESPACE}/healthcheckapp

docker build -t ${image}:${TAG} .

#docker push ${image}:${TAG}
#docker push ${image}:latest
rm root.tar
