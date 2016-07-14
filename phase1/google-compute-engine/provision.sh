#!/bin/bash -x

# Copyright 2016 The Kubernetes Authors All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

if ! /usr/bin/docker -v 2> /dev/null | grep -q "^Docker\ version\ 1\.10" ; then
  echo "Installing current version of Docker Engine 1.10"
  curl --silent --location  https://get.docker.com/builds/Linux/x86_64/docker-1.10.3  --output /usr/bin/docker
  chmod +x /usr/bin/docker
fi

systemd-run --unit=docker.service /usr/bin/docker daemon

/usr/bin/docker version

if ! [ -x /usr/bin/weave ] ; then
  echo "Installing current version of Weave Net"
  curl --silent --location https://git.io/weave --output /usr/bin/weave
  chmod +x /usr/bin/weave
  /usr/bin/weave setup
fi

/usr/bin/weave version

#/usr/bin/weave launch-router --init-peer-count 7

#/usr/bin/weave launch-proxy --rewrite-inspect

## Find nodes with `kube-weave` tag in an instance group

list_weave_peers_in_group() {
  ## There doesn't seem to be a native way to obtain instances with certain tags, so we use awk
  gcloud compute instance-groups list-instances $1 --uri --quiet \
    | xargs -n1 gcloud compute instances describe \
        --format='value(tags.items[], name, networkInterfaces[0].accessConfigs[0].natIP)' \
    | awk '$1 ~ /(^|\;)kube-weave($|\;).*/ && $2 ~ /^kube-.*$/ { print $2 }'
}

## This is very basic way of doing Weave Net peer discovery, one could potentially implement a pair of
## systemd units that write and watch an environment file and call `weave connect` when needed...
## However, the purpose of this script is only too illustrate the usage of Kubernetes Anywhere in GCE.

#/usr/bin/weave connect \
#  $(list_weave_peers_in_group kube-master-group) \
#  $(list_weave_peers_in_group kube-node-group)

#/usr/bin/weave expose -h $(hostname).weave.local

echo "MASTERS=$(list_weave_peers_in_group kube-master-group)" > /etc/ka2.env
echo "NODES=$(list_weave_peers_in_group kube-node-group)" >> /etc/ka2.env
