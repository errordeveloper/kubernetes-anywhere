#!/bin/sh
make crosscompile

IP1=`docker-machine ip vm1`
IP2=`docker-machine ip vm2`

docker-machine scp ./kube vm1:/tmp/kube
docker-machine scp ./kube vm2:/tmp/kube

if [ "$1" = "auto" ]; then
    docker-machine ssh vm1 sudo /tmp/kube init "${IP1},${IP2}" &
    pid="$$"
    docker-machine ssh vm2 sudo /tmp/kube join "${IP1},${IP2}"
    wait "$pid" 2>/dev/null
else
    echo
    echo "Log into vm1 with:"
    echo "    docker-machine ssh vm1"
    echo "Then run:"
    echo "    sudo /tmp/kube init ${IP1},${IP2}'"
    echo
    echo "To log into vm2, use:"
    echo "    docker-machine ssh vm2"
    echo
fi
