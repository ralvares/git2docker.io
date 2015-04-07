#!/bin/bash

docker pull cooltrick/git2docker
docker pull cooltrick/git2docker:start
docker pull busybox

mkdir /opt/git2docker
curl https://raw.githubusercontent.com/cooltrick/git2docker.io/master/gitreceive
curl https://github.com/cooltrick/git2docker.io/raw/master/receiver