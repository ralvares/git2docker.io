all:
	go get
	go build receiver.go 
	go build git2docker-cli.go 

install:
	mkdir /opt/git2docker
	cp -rf receiver /opt/git2docker
	cp -rf git2docker /opt/git2docker
	cp -rf git2docker-ssh /opt/git2docker
	cp -rf git2docker-cli /opt/git2docker
	cp -rf gitreceive	/opt/git2docker
	ln -s /opt/git2docker/git2docker /usr/bin/git2docker
	chmod +x /opt/git2docker/*
