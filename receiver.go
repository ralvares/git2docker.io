package main

import (
	//"fmt"
	"github.com/cooltrick/git2docker.io/build"
	"os"
)

var (
	userhome string = os.Getenv("HOME")
	username string = os.Getenv("USER")
	appname  string = os.Args[1]
	rev      string = os.Args[2]
	tmpdir   string = os.TempDir() + "/" + username + "_" + appname
)

func main() {

	build.BuildImage(appname, tmpdir, userhome, username, rev)

}
