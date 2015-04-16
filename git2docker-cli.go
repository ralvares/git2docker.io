package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

var (
	remove             = flag.Bool("remove", false, "Remove application")
	start              = flag.Bool("start", false, "Starting application")
	scale              = flag.Int("scale", 1, "-scale=X")
	ps                 = flag.Bool("ps", false, "Listing application")
	logs               = flag.Bool("logs", false, "Logs of application")
	port               = flag.Bool("port", false, "port of application")
	stop               = flag.Bool("stop", false, "Stop application")
	flagVersion        = flag.Bool("v", false, "Display version")
	name               = flag.String("name", "", "-name=NAME_OF_APPLICATION")
	userhome    string = os.Getenv("HOME")
	username    string = os.Getenv("USER")
	appname     string
)

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

func State(name string) bool {
	state := cmd("docker inspect --format '{{ .State.Running }}' " + name + " | grep -i true")
	if state == true {
		return true
	} else {
		return false
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

func CleanUP(name string) {
	if State("StorageApp_" + username + "_" + name) {
		Stop("StorageApp_" + username + "_" + name)
	}

	if ContainerExist("StorageApp_" + username + "_" + name) {
		RemoveContainer("StorageApp_" + username + "_" + name)
	}

	if State("App_" + username + "_" + name) {
		Stop("App_" + username + "_" + name)
	}

	if ContainerExist("App_" + username + "_" + name) {
		RemoveContainer("App_" + username + "_" + name)
	}

	if ImageExist(username + "/" + name + ":dockerfile") {
		RemoveImages(username + "/" + name + ":dockerfile")
	}

	if _, err := os.Stat(userhome + "/" + name); err == nil {
		os.RemoveAll(userhome + "/" + name)

	}
	fmt.Printf("| %-20s\t%10s |\n", name, "Deleted")
}

func main() {
	flag.Parse()

	if len(flag.Args()) < -1 {
		flag.Usage()
		return
	}

	if *flagVersion {
		fmt.Println("0.1")
		return
	}

	if *ps {
		if len(username) <= 0 {
			flag.Usage()
			return
		} else {
			List(userhome)
		}
	}

	if *logs {
		if len(*name) <= 0 {
			flag.Usage()
			return
		} else {
			Logs("App_" + username + "_" + *name)
			return
		}
	}

	if *port {
		if len(*name) <= 0 {
			flag.Usage()
			return
		} else {
			fmt.Println("Ports")
			return
		}
	}

	if *stop {
		if len(*name) <= 0 {
			flag.Usage()
			return
		} else {
			if _, err := os.Stat(userhome + "/" + *name); err == nil {
				if strings.HasPrefix(*name, ".") != true {
					if State(*name) {
						fmt.Printf("| %-20s\t%10s |\n", *name, "Already Stoped")
					} else {
						if Stop("App_" + username + "_" + *name) {
							fmt.Printf("| %-20s\t%10s |\n", *name, "Stoped")
						}
					}
				}
			} else {
				fmt.Println("App not Found.")
			}
		}
	}

	if *remove {
		if len(*name) <= 0 {
			flag.Usage()
			return
		} else {
			if _, err := os.Stat(userhome + "/" + *name); err == nil {
				CleanUP(*name)
			} else {
				fmt.Println("App not Found.")
			}
		}
	}

	if *start {
		if len(*name) <= 0 {
			flag.Usage()
			return
		} else {
			if _, err := os.Stat(userhome + "/" + *name); err == nil {

				if State(*name) {
					fmt.Printf("| %-20s\t%10s |\n", *name, "is Already UP")
				} else {
					if Start("App_" + username + "_" + *name) {
						fmt.Printf("| %-20s\t%10s |\n", *name, "Started")
					}
				}

			} else {
				fmt.Println("App not Found.")
			}
		}

	}
}
