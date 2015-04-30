#!/bin/bash

docker pull cooltrick/git2docker
docker pull cooltrick/git2docker:start
docker pull busybox
docker pull cooltrick/nginx-proxy

mkdir /opt/git2docker

curl https://raw.githubusercontent.com/cooltrick/git2docker.io/master/git2docker.io -o /opt/git2docker/git2docker.io
curl https://raw.githubusercontent.com/cooltrick/git2docker.io/master/sshcommand -o /opt/git2docker/sshcommand
ln -s /opt/git2docker/git2docker.io /opt/git2docker/receiver
ln -s /opt/git2docker/git2docker.io /opt/git2docker/git2docker-cli
ln -s /opt/git2docker/git2docker.io /opt/git2docker/gitreceive
ln -s /opt/git2docker/git2docker.io /usr/bin/git2docker
chmod +x /opt/git2docker/*