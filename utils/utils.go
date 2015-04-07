package utils

import (
	"fmt"
	"os"
	"os/exec"
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

	errtar := cmd("cd " + tmpdir + " && tar c . | docker run -i -a stdin --name=StorageApp_" + os.Getenv("USER") + "_" + name + " -v /app busybox /bin/sh -c 'tar -xC /app'")
	if errtar != true {
		fmt.Println("Error ---> Deploying Code...")
		return false
	} else {
		return true
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

func Run(name string, tmpdir string) {
	err := cmd("docker run -i -d -P --name=App_" + os.Getenv("USER") + "_" + name + " --volumes-from=StorageApp_" + os.Getenv("USER") + "_" + name + " cooltrick/git2docker:start '/start'")
	if err != true {
		fmt.Println("Error ---> Starting Code...")

	} else {
		fmt.Println(name + " Started")
		Ports(name)
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
	state := cmdout("docker inspect -f '{{range $p, $conf := .NetworkSettings.Ports}}{{(index $conf 0).HostPort}} {{end}}' App_" + os.Getenv("USER") + "_" + name)
	if state != true {
		fmt.Println("No Ports")
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
}

func GetCid(name string) string {
	return "App_" + os.Getenv("USER") + "_" + name
}
