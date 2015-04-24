all:
	go get
	go build

install:
	mkdir /opt/git2docker
	mkdir /opt/git2docker/databases/
	touch /opt/git2docker/databases/mysql:5.5
	touch /opt/git2docker/databases/mysql:5.6
	touch /opt/git2docker/databases/redis
	touch /opt/git2docker/databases/postgresql
	cp -rf git2docker.io /opt/git2docker
	cp -rf git2docker /opt/git2docker
	ln -s /opt/git2docker/git2docker.io /opt/git2docker/receiver
	ln -s /opt/git2docker/git2docker.io /opt/git2docker/git2docker-cli
	ln -s /opt/git2docker/git2docker.io /opt/git2docker/gitreceive
	ln -s /opt/git2docker/git2docker.io /usr/bin/git2docker
	chmod +x /opt/git2docker/*
