git2docker.io

Install Server:
OS: OpenSuse 13.2

You will need install docker.

zypper install docker
systemctl enable docker
systemctl start docker
systemctl enable sshd
systemctl start sshd

Disable Firewall.


curl https://raw.githubusercontent.com/cooltrick/git2docker.io/master/install.sh | sh

Create a User:

useradd -m user
gpasswd -a user docker
passwd user



Usage:

Add SSH key: 

ssh-keygen

cat ~/.ssh/id_rsa.pub | ssh user@X.X.X.X git2docker


Testing using git2docker.conf ( like heroku )

git clone https://github.com/heroku/node-js-sample

cd node-js-sample

git init

Create a git2docker.conf file:

echo state=build > git2docker.conf

git add --all
git commit -m "build"
git remote add git2docker user@X.X.X.X:node-js-sample

git push git2docker master

==============================================================
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


@linux:/tmp/node-js-sample> curl http://localhost:49153/
Hello World!
==============================================================

Testing using Dockerfile

mkdir apache-demo
cd apache-demo

echo "FROM httpd:2.4" > Dockerfile
echo "EXPOSE 80" >> Dockerfile

Create a git2docker.conf file:

echo state=dockerfile > git2docker.conf

git init
git add --all
git commit -m "build"
git remote add git2docker user@X.X.X.X:apache-demo

git push git2docker master


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


@linux:/tmp/node-js-sample> curl http://localhost:49154/
<html><body><h1>It works!</h1></body></html>


To remove the Containers

Just change the state flag to remove or delete and 

git add --all
git commit -m "delete"
git push git2docker master

==============================================================

How it Works:

You can deploy a Container using git2docker.conf file or using Dockerfile.


Using git2docker.conf


git2docker.conf Options:

===============================
state options:

build - Build the application using the source code.
build:logs - Build the application using the source code and start the container showing logs.
delete or remove - Stop and remove the Container
stop - Stop the Container
start - Start a stoped Container
start:logs - Start a stoped Container and show logs
logs - Show logs of a Started Container
dockerfile or Dockerfile - Force the git2docker to use a Dockerfile


ex: state=build

===============================
domain Option:

ex: domain=app.linux.site

===============================
pre-exec Option:

Option is used to execute a command before start the application

ex: pre-exec=bundle exec rake db:create db:migrate db:seed


===============================
git Option:

If you have your code at an external repository like github, git2docker will download the git and deploy the application.

ex: git=https://github.com/heroku/node-js-sample

===============================


git2docker.conf example

state=build
domain=app.domain.lnx
pre-exec=bundle exec rake db:create db:migrate db:seed

===============================

Provided your DNS is setup to forward foo.bar.com to the a host running nginx-proxy, the request will be routed to a container with the domain (git2docker.conf) env var set.

If you not set the domain var, the nginx-proxy will publish the domain using appname.username

Nginx proxy Simple Demo:

Create a systemd service:

cd /etc/systemd/system
vi nginx-proxy.service

[Unit]
Description=nginx-proxy
After=docker.service
Requires=docker.service

[Service]
TimeoutStartSec=0
ExecStartPre=-/usr/bin/docker kill nginx-proxy
ExecStartPre=-/usr/bin/docker rm nginx-proxy 
ExecStartPre=/usr/bin/docker pull jwilder/nginx-proxy  
ExecStart=/usr/bin/docker run -d --name=nginx-proxy -p 80:80 -v /var/run/docker.sock:/tmp/docker.sock jwilder/nginx-proxy

[Install]
WantedBy=multi-user.target

====================================================
Enable and Start service:

systemctl enable /etc/systemd/system/nginx-proxy.service
systemctl start nginx-proxy.service

Deploy any app seting domain option in git2docker.conf and try do access.

If you haven't a DNS server, You can add the domain at /etc/hosts


For a simple dns server to test:

Install dnsmasq 

zypper install dnsmasq 

If you want use the domain *.git2docker just run:

echo address="/.git2docker/192.168.100.187" >> /etc/dnsmasq.conf
systemctl enable dnsmasq
systemctl start dnsmasq

Use the new dns server to test you git2docker


nslookup test.git2docker
Server:              192.168.100.187
Address:      192.168.100.187#53

Name:  test.git2docker
Address: 192.168.100.187

SSH - Cli Option:

ex:

ssh user@X.X.X.X ps

Listing Containers.

@linux:/tmp/> ssh demo@192.168.100.187 ps
| apache-demo                    is Up |

Stoping Containers.

ssh user@X.X.X.X stop apache-demo

@linux:/tmp/> ssh demo@192.168.100.187 ps
| apache-demo                   Stoped |


Starting Containers.

ssh user@X.X.X.X start apache-demo
| apache-demo                  Started |


Deleting Containers.

ssh user@X.X.X.X delete apache-demo
Please type yes or no and then press enter:
yes
| apache-demo                  Deleted |


Linux Cli Option

mkdir ~/.git2docker

create a file ~/.git2docker/git2docker.conf:
user=user
host=X.X.X.X

cd ~/
curl https://github.com/cooltrick/git2docker.io/raw/master/git2docker-client/linux/git2docker
chmod +x git2docker

./git2docker -ps
| apache-demo                    is Up |

./git2docker -stop --name=apache-demo
| apache-demo                   Stoped |

./git2docker -start --name=apache-demo
| apache-demo                  Started |

./git2docker -logs --name=apache-demo
172.17.0.3 - - [16/Apr/2015:15:08:45 +0000] "GET / HTTP/1.1" 200 22698 "-" "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2272.118 Safari/537.36"
172.17.0.3 - - [16/Apr/2015:15:08:46 +0000] "GET / HTTP/1.1" 200 22700 "-" "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2272.118 Safari/537.36"
172.17.0.3 - - [16/Apr/2015:15:13:31 +0000] "GET / HTTP/1.1" 200 22702 "-" "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2272.118 Safari/537.36"
172.17.0.3 - - [16/Apr/2015:15:13:32 +0000] "GET / HTTP/1.1" 200 22698 "-" "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2272.118 Safari/537.36"
172.17.0.3 - - [16/Apr/2015:15:13:32 +0000] "GET / HTTP/1.1" 200 22702 "-" "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2272.118 Safari/537.36"
172.17.0.3 - - [16/Apr/2015:15:13:32 +0000] "GET / HTTP/1.1" 200 22698 "-" "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2272.118 Safari/537.36"
172.17.0.3 - - [16/Apr/2015:15:13:33 +0000] "GET / HTTP/1.1" 200 22697 "-" "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2272.118 Safari/537.36"

./git2docker -remove --name=apache-demo
Please type yes or no and then press enter:
yes
| apache-demo                  Deleted |