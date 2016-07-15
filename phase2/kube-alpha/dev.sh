#!/bin/sh
make crosscompile

IP1=`docker-machine ip vm1`
IP2=`docker-machine ip vm2`

docker-machine scp ./kube vm1:/usr/local/bin/kube
docker-machine scp ./kube vm2:/usr/local/bin/kube

docker-machine ssh vm1 /usr/local/bin/kube init $IP2
docker-machine ssh vm2 /usr/local/bin/kube join $IP1
