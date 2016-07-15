#!/bin/sh
docker-machine ssh vm1 'docker ps -a -q | xargs -r docker rm -f -v'
docker-machine ssh vm1 'rm -f /tmp/kube *.log'
docker-machine ssh vm2 'docker ps -a -q | xargs -r docker rm -f -v'
docker-machine ssh vm2 'rm -f /tmp/kube *.log'
