package build

import (
	"fmt"
	"github.com/cooltrick/cfg"
	"github.com/cooltrick/git2docker.io/utils"
	"log"
	"os"
	"os/exec"
	"strings"
)

type Git2Dockerconf struct {
	Domain     string
	State      string
	Preexec    string
	Git        string
	Database   bool
	Dockerfile bool
	Port       string
	Cache      bool
	login      string
	database   string
	password   string
	image      string
}

func (n *Git2Dockerconf) GetInfos(name string, tmpdir string, rev string) (string, string, string, string, bool, bool, string, bool, string, string, string, string) {

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

	if _, err := os.Stat(tmpdir + "_conf/Dockerfile"); err == nil {
		n.Dockerfile = true
	}

	if _, err := os.Stat(tmpdir + "_conf/git2docker.conf"); err != nil {
		File, errFile := os.Create(tmpdir + "_conf/git2docker.conf")
		defer File.Close()
		File.WriteString("\n")
		File.Sync()
		if errFile != nil {
			panic(errFile)
		}
	}

	myconf := make(map[string]string)
	err := cfg.Load(tmpdir+"_conf/git2docker.conf", myconf)
	if err != nil {
		log.Fatal(err)
	}

	for k, v := range myconf {
		if k == "state" {
			n.State = v
		}

		if k == "preexec" {
			n.Preexec = v
		}

		if k == "port" {
			n.Port = v
		}

		if k == "domain" {
			n.Domain = v
		}

		if k == "git" {
			n.Git = v
		}

		if k == "cache" {
			if v == "true" {
				n.Cache = true
			} else {
				n.Cache = false
			}
		}

	}

	if _, err := os.Stat(tmpdir + "_conf/git2docker_db.conf"); err == nil {
		myconf := make(map[string]string)
		err := cfg.Load(tmpdir+"_conf/git2docker_db.conf", myconf)
		if err != nil {
			log.Fatal(err)
		}

		for k, v := range myconf {

			if k == "user" {
				n.login = v
			}

			if k == "password" {
				n.password = v
			}

			if k == "database" {
				n.database = v
			}

			if k == "image" {
				n.image = v

			}
		}
	}
	//os.RemoveAll(os.TempDir() + "/" + os.Getenv("USER") + "_" + name + "_git2docker.conf")
	return n.Domain, n.State, n.Preexec, n.Git, n.Database, n.Dockerfile, n.Port, n.Cache, n.image, n.database, n.login, n.password
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

	if len(n.State) <= 0 || n.State == "delete" || n.State == "remove" || n.State == "stop" || n.State == "start" || n.State == "logs" || n.State == "start:logs" {
		n.Git = ""
	}

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

	if len(n.State) > 1 {
		n.Dockerfile = false
	}

	if _, err := os.Stat(tmpdir + "/git2docker_db.conf"); err == nil {
		if strings.HasSuffix(n.State, "rebuild") {
			if utils.State("DatabaseApp_" + os.Getenv("USER") + "_" + appname) {
				utils.Stop("DatabaseApp_" + os.Getenv("USER") + "_" + appname)
			}

			if utils.ContainerExist("DatabaseApp_" + os.Getenv("USER") + "_" + appname) {
				utils.RemoveContainer("DatabaseApp_" + os.Getenv("USER") + "_" + appname)
			}
		}

		if strings.HasPrefix(n.image, "redis") {
			if utils.State("DatabaseApp_" + os.Getenv("USER") + "_" + appname) {
				fmt.Println("Checking DataBase ...")
			} else {
				fmt.Println("Creating Database -> Redis ...")
				if utils.CreateDB(appname, n.image, "", "", "") {
					n.Database = true
				}
			}
		}

		if strings.HasPrefix(n.image, "mysql") {

			if utils.State("DatabaseApp_" + os.Getenv("USER") + "_" + appname) {
				fmt.Println("Checking DataBase ...")
			} else {
				fmt.Println("Creating Database - > Mysql ...")
				fmt.Println(os.Getenv("dbuser") + " " + os.Getenv("password") + " " + os.Getenv("dbname"))
				if utils.CreateDB(appname, n.image, n.login, n.password, n.database) {
					if _, err := os.Stat(tmpdir + "/git2docker.sql"); err == nil {
						var dbcheck bool
						fmt.Println("Waiting for " + n.database + " ...")
						fmt.Println("Press CTRL+C to Cancel")
						for dbcheck != true {
							dbcheck = utils.CMD("docker exec -i DatabaseApp_" + os.Getenv("USER") + "_" + appname + " mysqlshow -u" + n.login + " -p" + n.password + " " + n.database + " > /dev/null 2>&1 && true")
						}
						utils.CMD("docker exec -i DatabaseApp_" + os.Getenv("USER") + "_" + appname + " mysql -u" + n.login + " -p" + n.password + " " + n.database + " < " + tmpdir + "/git2docker.sql")
					}
					n.Database = true
					os.RemoveAll(tmpdir + "/git2docker.sql")
					os.RemoveAll(tmpdir + "/git2docker_db.conf")
				}
			}
		}
	} else {

		if utils.State("DatabaseApp_" + os.Getenv("USER") + "_" + appname) {
			utils.Stop("DatabaseApp_" + os.Getenv("USER") + "_" + appname)
		}

		if utils.ContainerExist("DatabaseApp_" + os.Getenv("USER") + "_" + appname) {
			utils.RemoveContainer("DatabaseApp_" + os.Getenv("USER") + "_" + appname)
		}
	}

	if n.Dockerfile {
		if utils.Dockerbuild(appname, tmpdir) {
			if n.Cache != true {
				fmt.Println("Cleaning Cache")
				os.RemoveAll(tmpdir)
			}
			utils.RunDockerbuild(appname, tmpdir, n.Domain)
		}
	}

	if n.State == "dockerfile" || n.State == "Dockerfile" || n.State == "dockerfile:rebuild" || n.State == "Dockerfile:rebuild" {
		n.Dockerfile = true
		if n.Dockerfile {
			if utils.Dockerbuild(appname, tmpdir) {
				if n.Cache != true {
					os.RemoveAll(tmpdir)
				}
				if n.Database {
					utils.RunDockerbuildwithDB(appname, tmpdir, n.Domain)
				} else {
					utils.RunDockerbuild(appname, tmpdir, n.Domain)
				}
			}
		}
	}

	if n.State == "build" || n.State == "build:rebuild" {
		os.RemoveAll(tmpdir + "/Dockerfile")
		os.RemoveAll(tmpdir + "/git2docker.conf")
		os.RemoveAll(tmpdir + "/git2docker_db.conf")
		if utils.CommitSource(appname, tmpdir) {
			if utils.Build(appname, tmpdir) {
				if n.Database {
					utils.RunwithDB(appname, tmpdir, n.Domain, n.Preexec)
				} else {
					utils.Run(appname, tmpdir, n.Domain, n.Preexec)
				}
			}
		}
	}

	if n.State == "build:logs" {
		os.RemoveAll(tmpdir + "/Dockerfile")
		os.RemoveAll(tmpdir + "/git2docker.conf")
		os.RemoveAll(tmpdir + "/git2docker_db.conf")
		if utils.CommitSource(appname, tmpdir) {
			if utils.Build(appname, tmpdir) {
				if n.Database {
					utils.RunwithDB(appname, tmpdir, n.Domain, n.Preexec)
					utils.Logs("App_" + username + "_" + appname)
				} else {
					utils.Run(appname, tmpdir, n.Domain, n.Preexec)
					utils.Logs("App_" + username + "_" + appname)
				}

			}
		}
	}

	if n.State == "logs" {
		if utils.State("App_" + username + "_" + appname) {
			utils.Logs("App_" + username + "_" + appname)
		} else {
			fmt.Println("APP don't exist...")
		}
	}

	if n.State == "remove" || n.State == "delete" {

		utils.CleanUP(appname)

		fmt.Println("App - " + appname + " - Removed")
	}

	if n.State == "env" {
		utils.GetEnv("App_" + username + "_" + appname)
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
			utils.Logs("App_" + username + "_" + appname)
		}
	}
}
