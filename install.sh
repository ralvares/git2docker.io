#!/bin/bash

docker pull cooltrick/git2docker
docker pull cooltrick/git2docker:start
docker pull busybox
docker pull cooltrick/nginx-proxy

mkdir /opt/git2docker
mkdir /opt/git2docker/databases/
touch /opt/git2docker/databases/mysql:5.5
touch /opt/git2docker/databases/mysql:5.6
touch /opt/git2docker/databases/redis
touch /opt/git2docker/databases/postgresql

curl https://raw.githubusercontent.com/cooltrick/git2docker.io/master/git2docker.io -o /opt/git2docker/git2docker.io
curl https://raw.githubusercontent.com/cooltrick/git2docker.io/master/git2docker -o /opt/git2docker/git2docker
ln -s /opt/git2docker/git2docker.io /opt/git2docker/receiver
ln -s /opt/git2docker/git2docker.io /opt/git2docker/git2docker-cli
ln -s /opt/git2docker/git2docker.io /opt/git2docker/gitreceive
ln -s /opt/git2docker/git2docker.io /usr/bin/git2docker
chmod +x /opt/git2docker/*