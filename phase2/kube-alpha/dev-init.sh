#!/bin/sh
export VIRTUALBOX_BOOT2DOCKER_URL="https://github.com/boot2docker/boot2docker/releases/download/v1.11.2/boot2docker.iso"
docker-machine create -d "${DOCKER_MACHINE_DRIVER:-virtualbox}" vm1
docker-machine create -d "${DOCKER_MACHINE_DRIVER:-virtualbox}" vm2
