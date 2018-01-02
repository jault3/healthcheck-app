#!/bin/bash -e
rm -f root/bin/healthcheckapp
go build -o root/bin/healthcheckapp ../main.go
rm -f root.tar
(cd root && fakeroot tar cvf ../root.tar .)
namespace=datica
image=registry-sbox05.datica.com/${namespace}/healthcheck-app
tag=2.0

docker build -t ${image}:${tag} .

docker push ${image}:${tag}
docker push ${image}:latest
