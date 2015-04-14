all:
	go build receiver.go 

install:
	cp -rf receiver /opt
	cp -rf git2docker /opt
	cp -rf git2docker-ssh /opt
	ln -s /opt/git2docker /usr/bin/git2docker
