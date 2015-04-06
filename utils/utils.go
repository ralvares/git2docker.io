package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
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
		fmt.Println("ERRO")
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

func Build(name string) bool {
	RemoveCid(name)
	errbuild := cmdout("docker run -i --cidfile=" + os.Getenv("HOME") + "/" + name + "/git2docker.cidfile " + os.Getenv("USER") + "/" + name + ":build /bin/bash -c '/build/builder'")
	if errbuild != true {
		fmt.Println("Error ---> Compiling Code...")
		RemoveContainer(GetCid(name))
		RemoveCid(name)
		return false
	} else {
		cmd("docker commit " + GetCid(name) + " " + os.Getenv("USER") + "/" + name + ":start")
		RemoveContainer(GetCid(name))
		RemoveCid(name)
		RemoveImages(os.Getenv("USER") + "/" + name + ":build")
		return true
	}
}

func CommitSource(name string, tmpdir string) bool {
	CleanUP(name)
	RemoveCid(name)
	errtar := cmd("cd " + tmpdir + " && tar c . | docker run -i -a stdin --cidfile=" + os.Getenv("HOME") + "/" + name + "/git2docker.cidfile build:image /bin/bash -c 'mkdir -p /app && tar -xC /app'")
	if errtar != true {
		fmt.Println("Error ---> Deploying Code...")
		RemoveContainer(GetCid(name))
		RemoveCid(name)
		return false
	} else {
		if cmd("docker commit " + GetCid(name) + " " + os.Getenv("USER") + "/" + name + ":build") {
			RemoveContainer(GetCid(name))
			return true
		}
	}
	return false

}

func RemoveContainer(name string) {
	//err := cmd("docker kill " + name + " && docker rm " + name)
	err := cmd("docker rm " + name)
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

func Stop(name string) {
	err := cmd("docker stop " + name)
	if err != true {
		fmt.Println("Error ---> Stoping Container Docker...")

	}
}

func Start(name string) {
	err := cmd("docker start " + name)
	if err != true {
		fmt.Println("Error ---> Startging Container Docker...")

	}
}

func Run(name string) {
	err := cmd("docker run -i -d -P --cidfile=" + os.Getenv("HOME") + "/" + name + "/git2docker.cidfile " + os.Getenv("USER") + "/" + name + ":start /start")
	if err != true {
		fmt.Println("Error ---> Starting Code...")

	}
}

func State(name string) bool {
	state := cmd("docker inspect --format '{{ .State.Running }}' " + name)
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

func Ports(name string) {
	state := cmd("docker inspect -f '{{range $p, $conf := .NetworkSettings.Ports}}{{(index $conf 0).HostPort}} {{end}}' " + name)
	if state != true {
		fmt.Println("No Ports")
	} else {
		fmt.Println(state)
	}
}

func CleanUP(name string) {
	if State(GetCid(name)) {
		Stop(GetCid(name))
	}

	if ContainerExist(GetCid(name)) {
		RemoveContainer(GetCid(name))
	}

	if ImageExist(os.Getenv("USER") + "/" + name + ":start") {
		RemoveImages(os.Getenv("USER") + "/" + name + ":start")
	}

	if ImageExist(os.Getenv("USER") + "/" + name + ":build") {
		RemoveImages(os.Getenv("USER") + "/" + name + ":build")
	}

	RemoveCid(name)

}

func CleanSource(name string) {

	if State(GetCid(name)) {
		Stop(GetCid(name))
	}

	if ContainerExist(GetCid(name)) {
		RemoveContainer(GetCid(name))
	}
	RemoveCid(name)
}

func GetCid(name string) string {
	content, err := ioutil.ReadFile(os.Getenv("HOME") + "/" + name + "/git2docker.cidfile")
	if err == nil {
		lines := strings.Split(string(content), "\n")
		return lines[0]
	} else {
		return "Error - GetCid"
	}
}

func RemoveCid(name string) {
	os.RemoveAll(os.Getenv("HOME") + "/" + name + "/git2docker.cidfile")
}
