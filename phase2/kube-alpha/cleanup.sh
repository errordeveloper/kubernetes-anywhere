#!/bin/sh
docker-machine ssh vm1 'docker ps -a -q | xargs -r docker rm -f'
docker-machine ssh vm2 'docker ps -a -q | xargs -r docker rm -f'
