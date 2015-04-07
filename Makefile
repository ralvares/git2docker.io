all:
	go build receiver.go 

install:
	cp -rf receiver /opt
