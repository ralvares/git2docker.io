all:
	go build -o git2docker.io main.go
	go build post-receive.go

install:
	cp -rf post-receive /opt
	ln -s /opt/post-receive /usr/share/git-core/templates/hooks/post-receive
