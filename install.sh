#!/bin/bash

docker pull cooltrick/git2docker
docker pull cooltrick/git2docker:start
docker pull busybox

mkdir /opt/git2docker
curl https://raw.githubusercontent.com/cooltrick/git2docker.io/master/gitreceive -o /opt/git2docker/gitreceive
curl https://raw.githubusercontent.com/cooltrick/git2docker.io/master/git2docker-ssh -o /opt/git2docker/git2docker-ssh
curl https://raw.githubusercontent.com/cooltrick/git2docker.io/master/git2docker -o /opt/git2docker/git2docker
curl https://github.com/cooltrick/git2docker.io/raw/master/receiver -o /opt/git2docker/receiver