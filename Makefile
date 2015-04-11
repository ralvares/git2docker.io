all:
	go build receiver.go 

install:
	cp -rf receiver /opt
	cp -rf git2docker /opt
	cp -rf git2docker-ssh /opt
