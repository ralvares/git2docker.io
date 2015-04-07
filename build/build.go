package build

import (
	"fmt"
	"github.com/cooltrick/cfg"
	"github.com/cooltrick/git2docker.io/utils"
	"log"
	"os"
	"os/exec"
)

func BuildImage(name string, tmpdir string, userhome string, username string, appname string) {
	if utils.VerifyAppName(name) {
		BuildAppGit(name, tmpdir, userhome, username, appname)
	} else {
		fmt.Println("Erro - VerifyAppName")
	}
}

func BuildAppGit(appname string, tmpdir string, userhome string, username string, rev string) {

	os.RemoveAll(tmpdir)

	GitCmd := exec.Command("git", "clone", userhome+"/"+appname, tmpdir)
	out, err := GitCmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(out))
		os.RemoveAll(tmpdir)
		panic(err)
	}

	os.Setenv("GIT_DIR", userhome+"/"+appname)
	os.Setenv("GIT_WORK_TREE", tmpdir)

	GitChk := exec.Command("git", "checkout", rev, "-f")
	out, errchk := GitChk.CombinedOutput()
	if errchk != nil {
		fmt.Println(string(out))
		os.RemoveAll(tmpdir)
		panic(errchk)
	}

	errChmod := os.Chmod(tmpdir, 0755)
	if errChmod != nil {
		os.RemoveAll(tmpdir)
		panic(errChmod)
	}

	if !utils.Createlock(tmpdir) {
		fmt.Println("Erro - LOCK FILE")
		os.RemoveAll(tmpdir)
		os.Exit(1)
	}

	//Building
	os.RemoveAll(tmpdir + "/.git")

	myconf := make(map[string]string)
	errconf := cfg.Load(tmpdir+"git2docker.conf", myconf)
	if errconf != nil {
		fmt.Printf("File git2docker.conf not Found")
		log.Fatal(err)
	}

	for k, v := range myconf {

		if k == "state" {

			if v == "build" {
				os.RemoveAll(tmpdir + "/git2docker.conf ")

				if utils.CommitSource(appname, tmpdir) {
					utils.Build(appname, tmpdir)
					utils.Run(appname, tmpdir)
				}
			}

			if v == "remove" || v == "delete" {

				utils.CleanUP(appname)

				fmt.Println("App - " + appname + " - Removed")
			}

			if v == "stop" {

				utils.Stop(utils.GetCid(appname))

				fmt.Println("App - " + appname + " - Stoped")
			}

			if v == "start" {
				if utils.State(utils.GetCid(appname)) {
					fmt.Println("Container -> UP")
				} else {
					utils.Start(utils.GetCid(appname))
					fmt.Println("App - " + appname + " - Started")
				}
			}

		}
	}
}
