all:
	go build receiver.go 
	go build git2docker-cli.go 

install:
	mkdir /opt/git2docker
	cp -rf receiver /opt/git2docker
	cp -rf git2docker /opt/git2docker
	cp -rf git2docker-ssh /opt/git2docker
	ln -s /opt/git2dockergit2docker /usr/bin/git2docker
