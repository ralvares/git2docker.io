# git2docker.io
# Read the Wiki
https://github.com/cooltrick/git2docker.io/wiki

git2docker.conf

Options:

state=build,delete ,stop,logs,build:logs,start and start:logs ( if logs flag is active, the git client will show the logs until Ctrl+C)


To Develop:

domain=domain.tlp
pre-exec=Command to execute before /start
git=http://login:pass@link/repo.git ( Auth just works with http )
database=mysql or pgsql or redis 