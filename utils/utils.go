package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

var (
	userhome string = os.Getenv("HOME")
	username string = os.Getenv("USER")
)

func Createlock(dirname string) bool {
	lockFile, err := os.Create(dirname + "/.lock")
	defer lockFile.Close()
	if err != nil {
		panic(err)
		return false
	} else {
		return true
	}

}

func VerifyLock(dirname string) bool {
	lockFile, err := os.Open(dirname + "/.lock")
	defer lockFile.Close()
	if err != nil {
		return false
	} else {
		return true
	}
}

func RemoveLock(dirname string) bool {
	err := os.RemoveAll(dirname + "/.lock")
	if err != nil {
		return false
	} else {
		return true
	}
}

func VerifyAppName(name string) bool {
	if len(name) > 0 {
		return true
	} else {
		return false
	}
}

func cmdout(command string) bool {
	cmd := exec.Command("bash", "-c", command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return false
		panic(err)
	} else {
		return true
	}
}

func cmd(command string) bool {
	cmd, err := exec.Command("bash", "-c", command+" >/dev/null").Output()
	if err != nil {
		fmt.Printf(string(cmd))
		return false
	} else {
		return true
	}
}

func CMD(command string) bool {
	cmd := exec.Command("bash", "-c", command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return false
		panic(err)
	} else {
		return true
	}
}

func Build(name string, tmpdir string) bool {

	if State("App_" + os.Getenv("USER") + "_" + name) {
		Stop("App_" + os.Getenv("USER") + "_" + name)
	}

	if ContainerExist("App_" + os.Getenv("USER") + "_" + name) {
		RemoveContainer("App_" + os.Getenv("USER") + "_" + name)
	}
	errbuild := cmdout("docker run -i --rm --name=App_" + os.Getenv("USER") + "_" + name + " --volumes-from StorageApp_" + os.Getenv("USER") + "_" + name + " cooltrick/git2docker /bin/bash -c '/build/builder'")
	if errbuild != true {
		fmt.Println("Error ---> Building Image...")
		if State("StorageApp_" + os.Getenv("USER") + "_" + name) {
			Stop("StorageApp_" + os.Getenv("USER") + "_" + name)
		}

		if ContainerExist("StorageApp_" + os.Getenv("USER") + "_" + name) {
			RemoveContainer("StorageApp_" + os.Getenv("USER") + "_" + name)
		}
		//RemoveContainer(GetCid(name))
		//RemoveCid(name)
		return false
	} else {
		//cmd("docker commit " + GetCid(name) + " " + os.Getenv("USER") + "/" + name)
		//RemoveContainer(GetCid(name))
		//RemoveCid(name)
		return true
	}
}

func CommitSource(name string, tmpdir string) bool {

	if State("StorageApp_" + os.Getenv("USER") + "_" + name) {
		Stop("StorageApp_" + os.Getenv("USER") + "_" + name)
	}

	if ContainerExist("StorageApp_" + os.Getenv("USER") + "_" + name) {
		RemoveContainer("StorageApp_" + os.Getenv("USER") + "_" + name)
	}

	if State("App_" + os.Getenv("USER") + "_" + name) {
		Stop("App_" + os.Getenv("USER") + "_" + name)
	}

	if ContainerExist("App_" + os.Getenv("USER") + "_" + name) {
		RemoveContainer("App_" + os.Getenv("USER") + "_" + name)
	}

	errtar := cmd("cd " + tmpdir + " && tar c . | docker run -i -a stdin --name=StorageApp_" + os.Getenv("USER") + "_" + name + " -v /app -v /tmp/cache:/cache busybox /bin/sh -c 'tar -xC /app'")
	if errtar != true {
		fmt.Println("Error ---> Deploying Code...")
		os.RemoveAll(tmpdir)
		return false
	} else {
		os.RemoveAll(tmpdir)
		return true
	}

}

func Dockerbuild(name string, tmpdir string) bool {

	if State("StorageApp_" + os.Getenv("USER") + "_" + name) {
		Stop("StorageApp_" + os.Getenv("USER") + "_" + name)
	}

	if ContainerExist("StorageApp_" + os.Getenv("USER") + "_" + name) {
		RemoveContainer("StorageApp_" + os.Getenv("USER") + "_" + name)
	}

	if State("App_" + os.Getenv("USER") + "_" + name) {
		Stop("App_" + os.Getenv("USER") + "_" + name)
	}

	if ContainerExist("App_" + os.Getenv("USER") + "_" + name) {
		RemoveContainer("App_" + os.Getenv("USER") + "_" + name)
	}

	if ImageExist(os.Getenv("USER") + "/" + name + ":dockerfile") {
		RemoveImages(os.Getenv("USER") + "/" + name + ":dockerfile")
	}

	errtar := cmdout("cd " + tmpdir + " && docker build --rm -t " + os.Getenv("USER") + "/" + name + ":dockerfile --force-rm .")
	if errtar != true {
		fmt.Println("Error ---> Docker Build...")
		os.RemoveAll(tmpdir)
		return false
	} else {
		return true
	}

}

func RunDockerbuild(name string, tmpdir string, domain string) {
	err := cmd("docker run -i -d -P --restart=always --name=App_" + os.Getenv("USER") + "_" + name + " -e VIRTUAL_HOST=" + domain + " " + os.Getenv("USER") + "/" + name + ":dockerfile")
	if err != true {
		fmt.Println("Error ---> Starting Docker...")
		RemoveContainer("App_" + os.Getenv("USER") + "_" + name)
	} else {
		fmt.Println(name + " Started")
		fmt.Println("Access: http://" + domain + " or Port: " + Ports(name))

	}
}

func RunDockerbuildwithDB(name string, tmpdir string, domain string) {
	err := cmd("docker run -i -d -P --restart=always --name=App_" + os.Getenv("USER") + "_" + name + " --link DatabaseApp_" + os.Getenv("USER") + "_" + name + ":database" + " -e VIRTUAL_HOST=" + domain + " " + os.Getenv("USER") + "/" + name + ":dockerfile")
	if err != true {
		fmt.Println("Error ---> Starting Docker...")
		RemoveContainer("App_" + os.Getenv("USER") + "_" + name)
	} else {
		fmt.Println(name + " Started")
		fmt.Println("")
		fmt.Println("Printing Database Informations:")
		fmt.Println("")
		GetEnv("App_" + os.Getenv("USER") + "_" + name)
		fmt.Println("")
		fmt.Println("Access: http://" + domain + " or Port: " + Ports(name))

	}
}

func RemoveContainer(name string) {
	//err := cmd("docker kill " + name + " && docker rm " + name)
	err := cmd("docker rm -f -v " + name)

	if err != true {
		fmt.Println("Error ---> Deleting Container Docker...")

	}
}

func RemoveImages(name string) {
	err := cmd("docker rmi " + name)
	if err != true {
		fmt.Println("Error ---> Deleting Image Docker...")

	}
}

func Stop(name string) bool {
	err := cmd("docker stop " + name)
	if err != true {
		fmt.Println("Error ---> Stoping Container Docker...")
		return false
	} else {
		return true
	}
}

func Start(name string) bool {
	err := cmd("docker start " + name)
	if err != true {
		fmt.Println("Error ---> Starting Container Docker...")
		return false
	} else {
		return true
	}
}

func Logs(name string) {
	err := cmdout("docker logs --tail=1000 " + name)
	if err != true {
		fmt.Println("Error ---> Showing Logs Docker...")

	}
}

func Run(name string, tmpdir string, domain string, preexec string) {
	if len(preexec) < 0 {
		preexec = "echo OK"
	}
	errPre := cmd("docker run -i --rm --name=App_" + os.Getenv("USER") + "_" + name + " --volumes-from=StorageApp_" + os.Getenv("USER") + "_" + name + " cooltrick/git2docker:start /bin/bash -c '/preexec " + preexec + "'")
	if errPre != true {
		fmt.Println("Error ---> Starting Pre-Exec...")

	} else {
		if State("App_" + os.Getenv("USER") + "_" + name) {
			Stop("App_" + os.Getenv("USER") + "_" + name)
		}

		if ContainerExist("App_" + os.Getenv("USER") + "_" + name) {
			RemoveContainer("App_" + os.Getenv("USER") + "_" + name)
		}
		err := cmd("docker run -i -d -P --restart=always --name=App_" + os.Getenv("USER") + "_" + name + " --volumes-from=StorageApp_" + os.Getenv("USER") + "_" + name + " -e VIRTUAL_HOST=" + domain + " cooltrick/git2docker:start /bin/bash -c '/start'")
		if err != true {
			fmt.Println("Error ---> Starting Code...")
			RemoveContainer("App_" + os.Getenv("USER") + "_" + name)
		} else {
			fmt.Println(name + " Started")
			fmt.Println("")
			fmt.Println("Printing Database Informations:")
			fmt.Println("")
			GetEnv("App_" + os.Getenv("USER") + "_" + name)
			fmt.Println("")
			fmt.Println("Access: http://" + domain + " or Port: " + Ports(name))

		}
	}
}

func RunwithDB(name string, tmpdir string, domain string, preexec string) {
	if len(preexec) < 0 {
		preexec = "echo OK"
	}
	errPre := cmd("docker run -i --rm --name=App_" + os.Getenv("USER") + "_" + name + " --volumes-from=StorageApp_" + os.Getenv("USER") + "_" + name + " cooltrick/git2docker:start /bin/bash -c '/preexec " + preexec + "'")
	if errPre != true {
		fmt.Println("Error ---> Starting Pre-Exec...")

	} else {
		if State("App_" + os.Getenv("USER") + "_" + name) {
			Stop("App_" + os.Getenv("USER") + "_" + name)
		}

		if ContainerExist("App_" + os.Getenv("USER") + "_" + name) {
			RemoveContainer("App_" + os.Getenv("USER") + "_" + name)
		}
		err := cmd("docker run -i -d -P --restart=always --name=App_" + os.Getenv("USER") + "_" + name + " --link DatabaseApp_" + os.Getenv("USER") + "_" + name + ":database" + " --volumes-from=StorageApp_" + os.Getenv("USER") + "_" + name + " -e VIRTUAL_HOST=" + domain + " cooltrick/git2docker:start /bin/bash -c '/start'")
		if err != true {
			fmt.Println("Error ---> Starting Code...")

		} else {
			fmt.Println(name + " Started")
			fmt.Println("Access: http://" + domain + " or Port: " + Ports(name))
		}
	}
}

func CreateDB(name string, image string, login string, pass string, dbname string) bool {
	if State("DatabaseApp_" + os.Getenv("USER") + "_" + name) {
		Stop("DatabaseApp_" + os.Getenv("USER") + "_" + name)
	}

	if ContainerExist("DatabaseApp_" + os.Getenv("USER") + "_" + name) {
		RemoveContainer("DatabaseApp_" + os.Getenv("USER") + "_" + name)
	}

	if strings.HasPrefix(image, "redis") {
		err := cmd("docker run -i -d --restart=always --name=DatabaseApp_" + os.Getenv("USER") + "_" + name + " " + image)
		if err != true {
			fmt.Println("Error ---> Starting Database...")
			RemoveContainer("DatabaseApp_" + os.Getenv("USER") + "_" + name)
			return false
		} else {
			fmt.Println("DataBase Started")
			return true
		}
	}
	if strings.HasPrefix(image, "mysql") {
		err := cmd("docker run -i -d --restart=always --name=DatabaseApp_" + os.Getenv("USER") + "_" + name + " -e MYSQL_ROOT_PASSWORD=" + pass + " -e MYSQL_DATABASE=" + dbname + " -e MYSQL_USER=" + login + " -e MYSQL_PASSWORD=" + pass + " " + image)
		if err != true {
			fmt.Println("Error ---> Starting Database...")
			RemoveContainer("DatabaseApp_" + os.Getenv("USER") + "_" + name)
			return false
			os.Exit(1)
		} else {
			fmt.Println("DataBase Started")
			return true
		}
	}

	return false
}

/*func State(name string) bool {
	state := cmd("docker inspect --format '{{ .State.Running }}' " + name)
	if state == true {
		return true
	} else {
		return false
	}
}*/

func State(name string) bool {
	state := cmd("docker inspect --format '{{ .State.Running }}' " + name + " | grep -i true")
	if state == true {
		return true
	} else {
		return false
	}
}

func ImageExist(name string) bool {
	state := cmd("docker inspect " + name)
	if state != true {
		return false
	} else {
		return true
	}
}

func ContainerExist(name string) bool {
	state := cmd("docker inspect " + name)
	if state != true {
		return false
	} else {
		return true
	}
}

func Ports(name string) string {
	state := exec.Command("bash", "-c", "docker inspect -f '{{range $p, $conf := .NetworkSettings.Ports}}{{(index $conf 0).HostPort}} {{end}}' App_"+os.Getenv("USER")+"_"+name)
	out, err := state.CombinedOutput()
	if err != nil {
		panic(err)
	} else {
		return string(out)
	}
}

func GetEnv(name string) {
	err := cmdout("docker exec " + name + " /bin/bash -c export | grep -i database | awk '{print $3}'")
	if err != true {
		fmt.Println("Error ---> Container is Down")
	}
}

func List(userhome string) {
	files, _ := ioutil.ReadDir(userhome)
	for _, f := range files {
		if f.IsDir() {
			if strings.HasPrefix(f.Name(), ".") != true {
				if State("App_"+username+"_"+f.Name()) != true {
					fmt.Printf("| %-20s\t%10s |\n", f.Name(), "is Down")
				} else {
					fmt.Printf("| %-20s\t%10s |\n", f.Name(), "is Up")
				}
			}
		}
	}
}

func CleanUP(name string) {
	if State("StorageApp_" + os.Getenv("USER") + "_" + name) {
		Stop("StorageApp_" + os.Getenv("USER") + "_" + name)
	}

	if ContainerExist("StorageApp_" + os.Getenv("USER") + "_" + name) {
		RemoveContainer("StorageApp_" + os.Getenv("USER") + "_" + name)
	}

	if State("App_" + os.Getenv("USER") + "_" + name) {
		Stop("App_" + os.Getenv("USER") + "_" + name)
	}

	if ContainerExist("App_" + os.Getenv("USER") + "_" + name) {
		RemoveContainer("App_" + os.Getenv("USER") + "_" + name)
	}

	if ImageExist(os.Getenv("USER") + "/" + name + ":dockerfile") {
		RemoveImages(os.Getenv("USER") + "/" + name + ":dockerfile")
	}

	if State("DatabaseApp_" + os.Getenv("USER") + "_" + name) {
		fmt.Println("Stoping Database")
		Stop("DatabaseApp_" + os.Getenv("USER") + "_" + name)
	}

	if ContainerExist("DatabaseApp_" + os.Getenv("USER") + "_" + name) {
		fmt.Println("Deleting Database")
		RemoveContainer("DatabaseApp_" + os.Getenv("USER") + "_" + name)
	}

	if _, err := os.Stat(os.Getenv("HOME") + "/" + name); err == nil {
		File, errFile := os.Create(os.Getenv("HOME") + "/" + name + "/.remove")
		defer File.Close()
		if errFile != nil {
			panic(errFile)
		}
	}
}

func GetCid(name string) string {
	return "App_" + os.Getenv("USER") + "_" + name
}
