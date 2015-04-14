e # git2docker.io


To Install:

Test made using opensuse and centos.

curl https://raw.githubusercontent.com/cooltrick/git2docker.io/master/install.sh | sh


Add SSH key: 

cat ~/.ssh/id_rsa.pub | ssh LOGIN@X.X.X.X git2docker


git2docker.conf file

===============================
state Options:

build

build:logs

delete or remove

stop

start

start:logs

logs

dockerfile or Dockerfile


===============================
domain Option:

===============================
pre-exec Option:

===============================
git Option:



===============================
To Develop:

database=mysql or pgsql or redis