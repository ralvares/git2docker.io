package build

import (
	"fmt"
	"github.com/cooltrick/cfg"
	"github.com/cooltrick/git2docker.io/utils"
	"log"
	"os"
	"os/exec"
)

type Git2Dockerconf struct {
	Domain     string
	State      string
	Preexec    string
	Git        string
	Database   string
	Dockerfile bool
}

func (n *Git2Dockerconf) GetInfos(name string, tmpdir string, rev string) (string, string, string, string, string, bool) {

	os.RemoveAll(tmpdir)
	os.RemoveAll(tmpdir + "_conf")

	GitCmd := exec.Command("git", "clone", os.Getenv("HOME")+"/"+name, tmpdir+"_conf")
	out, errGit := GitCmd.CombinedOutput()
	if errGit != nil {
		fmt.Println(string(out))
		os.RemoveAll(tmpdir + "_conf")
		panic(errGit)
	}

	os.Setenv("GIT_DIR", os.Getenv("HOME")+"/"+name)
	os.Setenv("GIT_WORK_TREE", tmpdir+"_conf")

	GitChk := exec.Command("git", "checkout", rev, "-f")
	out, errchk := GitChk.CombinedOutput()
	if errchk != nil {
		fmt.Println(string(out))
		os.RemoveAll(tmpdir + "_conf")
		panic(errchk)
	}

	errChmod := os.Chmod(tmpdir+"_conf", 0755)
	if errChmod != nil {
		os.RemoveAll(tmpdir + "_conf")
		panic(errChmod)
	}

	myconf := make(map[string]string)
	err := cfg.Load(tmpdir+"_conf/git2docker.conf", myconf)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := os.Stat(tmpdir + "_conf/Dockerfile"); err == nil {
		n.Dockerfile = true
	}

	for k, v := range myconf {
		if k == "state" {
			n.State = v
		}
		if k == "preexec" {
			n.Preexec = v
		}

		if k == "domain" {
			n.Domain = v
		}

		if k == "database" {
			n.Database = v
		}
		if k == "git" {
			if n.Dockerfile {
				n.Git = ""
			} else {
				n.Git = v
			}
		}

	}
	//os.RemoveAll(os.TempDir() + "/" + os.Getenv("USER") + "_" + name + "_git2docker.conf")
	return n.Domain, n.State, n.Preexec, n.Git, n.Database, n.Dockerfile
}

func BuildImage(name string, tmpdir string, userhome string, username string, rev string) {
	if utils.VerifyAppName(name) {
		BuildAppGit(name, tmpdir, userhome, username, rev)
	} else {
		fmt.Println("Erro - VerifyAppName")
	}
}

func BuildAppGit(appname string, tmpdir string, userhome string, username string, rev string) {

	n := Git2Dockerconf{}
	n.GetInfos(appname, tmpdir, rev)

	if len(n.Git) <= 0 {

		MVCmd := exec.Command("mv", tmpdir+"_conf", tmpdir)
		out, err := MVCmd.CombinedOutput()
		if err != nil {
			fmt.Println(string(out))
			os.RemoveAll(tmpdir + "_conf")
			panic(err)
		}

	} else {

		os.Setenv("GIT_DIR", os.Getenv("HOME")+"/"+appname)
		os.Setenv("GIT_WORK_TREE", tmpdir)

		GitCmd := exec.Command("git", "clone", n.Git, tmpdir)
		out, err := GitCmd.CombinedOutput()
		if err != nil {
			fmt.Println(string(out))
			os.RemoveAll(tmpdir)
			panic(err)
		}

		errChmod := os.Chmod(tmpdir, 0755)
		if errChmod != nil {
			os.RemoveAll(tmpdir)
			panic(errChmod)
		}
	}

	//Building
	os.RemoveAll(tmpdir + "/.git")

	if len(n.Domain) <= 0 {
		n.Domain = appname + "." + os.Getenv("USER")
	}

	if len(n.Preexec) <= 0 {
		n.Preexec = "true"
	}

	if n.State == "build" {

		if utils.CommitSource(appname, tmpdir) {
			utils.Build(appname, tmpdir)
			utils.Run(appname, tmpdir, n.Domain, n.Preexec)
		}
	}

	if n.State == "Dockerfile" || n.State == "dockerfile" {
		if n.Dockerfile {
			fmt.Println("Dockerfile exist")
			if utils.Dockerbuild(appname, tmpdir) {
				utils.RunDockerbuild(appname, tmpdir, n.Domain)
			}
		} else {
			fmt.Println("Dockerfile don't exist")
		}
	}

	if n.State == "build:logs" {
		if utils.CommitSource(appname, tmpdir) {
			utils.Build(appname, tmpdir)
			utils.Run(appname, tmpdir, n.Domain, n.Preexec)
			utils.Logs(appname)
		}
	}

	if n.State == "logs" {
		if utils.State("App_" + username + "_" + appname) {
			utils.Logs(appname)
		} else {
			fmt.Println("APP don't exist...")
		}
	}

	if n.State == "remove" || n.State == "delete" {

		utils.CleanUP(appname)

		fmt.Println("App - " + appname + " - Removed")
	}

	if n.State == "stop" {

		utils.Stop(utils.GetCid(appname))

		fmt.Println("App - " + appname + " - Stoped")
	}

	if n.State == "start" {
		if utils.State(utils.GetCid(appname)) {
			fmt.Println("Container -> UP")
		} else {
			utils.Start(utils.GetCid(appname))
			fmt.Println("App - " + appname + " - Started")
		}
	}

	if n.State == "start:logs" {
		if utils.State(utils.GetCid(appname)) {
			fmt.Println("Container -> UP")
		} else {
			utils.Start(utils.GetCid(appname))
			fmt.Println("App - " + appname + " - Started")
			utils.Logs(appname)
		}
	}

}
