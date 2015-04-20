###﻿Git2Docker – Server


- Install

OS: Opensuse 13.2

```
zypper install docker
systemctl enable docker
systemctl start docker
systemctl enable sshd
systemctl start sshd
```
_Disable the Firewall_.

- Installing - Git2Docker

```
curl https://raw.githubusercontent.com/cooltrick/git2docker.io/master/install.sh | sh
```

- Installing using the Source.

```
zypper install go

git clone https://github.com/cooltrick/git2docker.io.git
cd git2docker.io
make
sudo make install

docker pull cooltrick/git2docker
docker pull cooltrick/git2docker:start
docker pull busybox
docker pull cooltrick/nginx-proxy
```

- Creating the user:

```
useradd -m user
gpasswd -a user docker
passwd user
```
###﻿Git2Docker – Client

```
mkdir ~/.git2docker
ssh-keygen -t dsa
cp -rf ~/.ssh/id_rsa* ~/.git2docker/
cat ~/.git2docker/id_rsa.pub | ssh user@192.168.100.56 git2docker
```

 - Downloading and installing the g2docker cli.

```
curl https://github.com/cooltrick/git2docker.io/raw/master/git2docker-client/linux/g2docker
chmod +x g2docker

cat >  ~/.git2docker/git2docker.conf <<EOF
user=user
host=192.168.100.56
EOF

```

- Installing g2docker cli using the Source.

```
zypper install go

git clone https://github.com/cooltrick/git2docker.io.git
cd git2docker.io/git2docker-client/linux
make
chmod +x g2docker
cp -rf g2docker /usr/local/bin
```


###About git2docker.conf

The file git2docker.conf is necessary, the server read the file and get needed options to work.

**git2docker.conf**:

state options:

- build - Build the application using the source code.
- build:logs - Build the application using the source code and start the container showing logs.
- delete or remove - Stop and remove the Container
- stop - Stop the Container
- start - Start a stoped Container
- start:logs - Start a stoped Container and show logs
- logs - Show logs of a Started Container
- dockerfile or Dockerfile - Force the git2docker to use a Dockerfile


example:

```
state=build
```

domain option:
> Provided your DNS is setup to forward your domain to the a host running nginx-proxy, the request will be routed to a container with the domain setted in git2docker.conf.


example:
```
domain=app.linux.site
```

pre-exec option:

>Option used to execute a command before start the application.

example:
```
pre-exec=bundle exec rake db:create db:migrate db:seed
```

git Option:

> If you have your code at an external repository like github, git2docker will download and deploy the application.

example:
```
git=https://github.com/heroku/node-js-sample
```


Exemple of git2docker.conf

```
state=build
domain=app.domain.lnx
pre-exec=bundle exec rake db:create db:migrate db:seed
```

###Deploy:

We will deploy a Node.js example application which starts a minimal HTTP server.
```
git clone https://github.com/heroku/node-js-sample
cd node-js-sample
git init

echo state=build > git2docker.conf
git add --all
git commit -m "build"
git remote add git2docker user@192.168.100.56:node-js-sample

git push git2docker master
```

```
@linux:/tmp/node-js-sample> git push git2docker master
Counting objects: 391, done.
Delta compression using up to 4 threads.
Compressing objects: 100% (316/316), done.
Writing objects: 100% (391/391), 214.55 KiB | 0 bytes/s, done.
Total 391 (delta 46), reused 387 (delta 45)
remote: =======> Working - node-js-sample
remote:
-----> Using u1000 to run an application
-----> Node.js app detected

-----> Reading application state
       package.json...
       build directory...
       cache directory...
       environment variables...

       Node engine:         0.12.x
       Npm engine:          unspecified
       Start mechanism:     npm start
       node_modules source: npm-shrinkwrap.json
       node_modules cached: false

       NPM_CONFIG_PRODUCTION=true
       NODE_MODULES_CACHE=true

-----> Installing binaries
       Resolving node version 0.12.x via semver.io...
       Downloading and installing node 0.12.2...
       Using default npm version: 2.7.4

-----> Building dependencies
       Installing node modules
       express@4.12.3 node_modules/express
       ├── merge-descriptors@1.0.0
       ├── utils-merge@1.0.0
       ├── cookie-signature@1.0.6
       ├── methods@1.1.1
       ├── fresh@0.2.4
       ├── escape-html@1.0.1
       ├── cookie@0.1.2
       ├── range-parser@1.0.2
       ├── content-type@1.0.1
       ├── finalhandler@0.3.4
       ├── vary@1.0.0
       ├── parseurl@1.3.0
       ├── serve-static@1.9.2
       ├── content-disposition@0.5.0
       ├── path-to-regexp@0.1.3
       ├── depd@1.0.1
       ├── qs@2.4.1
       ├── etag@1.5.1 (crc@3.2.1)
       ├── on-finished@2.2.0 (ee-first@1.1.0)
       ├── debug@2.1.3 (ms@0.7.0)
       ├── proxy-addr@1.0.7 (forwarded@0.1.0, ipaddr.js@0.1.9)
       ├── send@0.12.2 (destroy@1.0.3, ms@0.7.0, mime@1.3.4)
       ├── accepts@1.2.5 (negotiator@0.5.1, mime-types@2.0.10)
       └── type-is@1.6.1 (media-typer@0.3.0, mime-types@2.0.10)

-----> Checking startup method
       No Procfile; Adding 'web: npm start' to new Procfile

-----> Finalizing build
       Creating runtime environment
       Exporting binary paths
       Cleaning npm artifacts
       Cleaning previous cache
       Caching results for future builds

-----> Build succeeded!

       node-js-sample@0.2.0 /tmp/build
       └── express@4.12.3

-----> Discovering process types
       Procfile declares types -> web
remote: node-js-sample Started
remote: 49153
To demo@localhost:node-js-sample
* [new branch]      master -> master



@linux:/tmp/node-js-sample> curl http://192.168.100.56:49153/
Hello World!
```

###Deploy - Using a Dockerfile

We will deploy a Apache(httpd) example application.

```
mkdir apache-demo
cd apache-demo

echo "FROM httpd:2.4" > Dockerfile
echo "EXPOSE 80" >> Dockerfile

echo state=dockerfile > git2docker.conf

git add --all
git commit -m "build"
git remote add git2docker user@192.168.100.56:apache-demo

git push git2docker master
```
```
@linux:/tmp/apache-demo> git push git2docker master
Counting objects: 4, done.
Delta compression using up to 4 threads.
Compressing objects: 100% (2/2), done.
Writing objects: 100% (4/4), 303 bytes | 0 bytes/s, done.
Total 4 (delta 0), reused 0 (delta 0)
remote: =======> Working - apache-demo
remote:
remote: Sending build context to Docker daemon 3.072 kB
remote: Sending build context to Docker daemon
remote: Step 0 : FROM httpd:2.4
remote:  ---> 4ea677a2d898
remote: Step 1 : EXPOSE 80
remote:  ---> Running in 2fe3a7300cdf
remote:  ---> 11d40bb1d4fe
remote: Removing intermediate container 2fe3a7300cdf
remote: Successfully built 11d40bb1d4fe
remote: apache-demo Started
remote: 49154
To demo@localhost:apache-demo
 * [new branch]      master -> master


@linux:/tmp/node-js-sample> curl http://192.168.100.56:49154/
It works!
```

###Manage Containers - Git Client

- Deleting:

```
echo state=remove > git2docker.conf
git add --all
git commit -m "build"
git push git2docker master

```

- Stoping:

```
echo state=stop > git2docker.conf
git add --all
git commit -m "build"
git push git2docker master

```

- Starting:

```
echo state=start > git2docker.conf
git add --all
git commit -m "build"
git push git2docker master

```

- Logs:

```
echo state=logs > git2docker.conf
git add --all
git commit -m "build"
git push git2docker master

```

###Manage Containers - Git2Docker-CLI Client

- Listing:

```

./g2docker -ps
| apache-demo                    is Up |
```

- Stopting:

```

./g2docker -stop --name=apache-demo
| apache-demo                   Stoped |
```

- Starting:

```

./g2docker -start --name=apache-demo
| apache-demo                  Started |
```

- Logs:

```

./g2docker -logs --name=apache-demo
172.17.0.3 - - [16/Apr/2015:15:08:45 +0000] "GET / HTTP/1.1" 200 22698 "-" "
172.17.0.3 - - [16/Apr/2015:15:08:45 +0000] "GET / HTTP/1.1" 200 22698 "-" "
172.17.0.3 - - [16/Apr/2015:15:08:45 +0000] "GET / HTTP/1.1" 200 22698 "-" "
172.17.0.3 - - [16/Apr/2015:15:08:45 +0000] "GET / HTTP/1.1" 200 22698 "-" "
172.17.0.3 - - [16/Apr/2015:15:08:45 +0000] "GET / HTTP/1.1" 200 22698 "-" "
172.17.0.3 - - [16/Apr/2015:15:08:45 +0000] "GET / HTTP/1.1" 200 22698 "-" "
```

- Deleting:

```
./g2docker -remove --name=apache-demo
Please type yes or no and then press enter: yes
| apache-demo                  Deleted |
```

####Nginx proxy:

>If you not set the domain option, the nginx-proxy will publish the domain using appname.username

- Create a systemd service:

```
cd /etc/systemd/system

cat >  nginx-proxy.service  <<EOF
[Unit]
Description=nginx-proxy
After=docker.service
Requires=docker.service

[Service]
TimeoutStartSec=0
ExecStartPre=-/usr/bin/docker kill nginx-proxy
ExecStartPre=-/usr/bin/docker rm nginx-proxy
ExecStartPre=/usr/bin/docker pull cooltrick/nginx-proxy
ExecStart=/usr/bin/docker run -d --name=nginx-proxy -p 80:80 -v /var/run/docker.sock:/tmp/docker.sock cooltrick/nginx-proxy

[Install]
WantedBy=multi-user.target
EOF
```

 - Enable and Start service:

```
systemctl enable /etc/systemd/system/nginx-proxy.service
systemctl start nginx-proxy.service
```

Deploy any Container setting domain option in git2docker.conf and test.

>If you haven't a DNS server, You can add the domain at /etc/hosts.

Example of /etc/hosts file:
```
192.168.100.56	apache-demo.git2docker
192.168.100.56	nodejs.git2docker
```
