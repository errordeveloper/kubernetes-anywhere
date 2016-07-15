#!/bin/sh
make crosscompile

IP1=`docker-machine ip vm1`
IP2=`docker-machine ip vm2`

docker-machine scp ./kube vm1:/tmp/kube
docker-machine scp ./kube vm2:/tmp/kube

docker-machine ssh vm1 sudo /tmp/kube init $IP2
docker-machine ssh vm2 sudo /tmp/kube join $IP1
