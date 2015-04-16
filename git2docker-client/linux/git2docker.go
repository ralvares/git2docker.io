package main

import (
	"flag"
	"fmt"
	"github.com/cooltrick/cfg"
	"github.com/hypersleep/easyssh"
	"log"
	"os"
)

var (
	remove      = flag.Bool("remove", false, "Remove application")
	start       = flag.Bool("start", false, "Starting application")
	scale       = flag.Int("scale", 1, "-scale=X")
	ps          = flag.Bool("ps", false, "Listing application")
	logs        = flag.Bool("logs", false, "Logs of application")
	stop        = flag.Bool("stop", false, "Stop application")
	flagVersion = flag.Bool("v", false, "Display version")
	name        = flag.String("name", "", "-name=NAME_OF_APPLICATION")
	user        string
	host        string
)

func posString(slice []string, element string) int {
	for index, elem := range slice {
		if elem == element {
			return index
		}
	}
	return -1
}

func containsString(slice []string, element string) bool {
	return !(posString(slice, element) == -1)
}

func askForConfirmation() bool {
	fmt.Println("Please type yes or no and then press enter:")
	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		log.Fatal(err)
	}
	okayResponses := []string{"y", "Y", "yes", "Yes", "YES"}
	nokayResponses := []string{"n", "N", "no", "No", "NO"}
	if containsString(okayResponses, response) {
		return true
	} else if containsString(nokayResponses, response) {
		return false
	} else {
		fmt.Println("Please type yes or no and then press enter:")
		return askForConfirmation()
	}
}

func main() {

	myconf := make(map[string]string)
	err := cfg.Load(os.Getenv("HOME")+"/.git2docker/git2docker.conf", myconf)
	if err != nil {
		log.Fatal(err)
	}

	for k, v := range myconf {
		if k == "user" {
			user = v
		}
		if k == "host" {
			host = v
		}
	}

	ssh := &easyssh.MakeConfig{
		User:   user,
		Server: host,
		Key:    "/.git2docker/id_rsa",
		Port:   "22",
	}

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
		response, err := ssh.Run("ps")
		// Handle errors
		if err != nil {
			panic("Can't run remote command: " + err.Error())
		} else {
			fmt.Println(response)
		}

	}

	if *remove {
		if len(*name) <= 0 {
			flag.Usage()
			return
		} else {
			if askForConfirmation() {
				response, err := ssh.Run("remove " + *name)
				// Handle errors
				if err != nil {
					panic("Can't run remote command: " + err.Error())
				} else {
					fmt.Println(response)
				}
			}
		}

	}

	if *start {
		if len(*name) <= 0 {
			flag.Usage()
			return
		} else {
			response, err := ssh.Run("start " + *name)
			// Handle errors
			if err != nil {
				panic("Can't run remote command: " + err.Error())
			} else {
				fmt.Println(response)
			}
		}

	}

	if *stop {
		if len(*name) <= 0 {
			flag.Usage()
			return
		} else {
			response, err := ssh.Run("stop " + *name)
			// Handle errors
			if err != nil {
				panic("Can't run remote command: " + err.Error())
			} else {
				fmt.Println(response)
			}
		}

	}

	if *logs {
		if len(*name) <= 0 {
			flag.Usage()
			return
		} else {
			response, err := ssh.Run("logs " + *name)
			// Handle errors
			if err != nil {
				panic("Can't run remote command: " + err.Error())
			} else {
				fmt.Println(response)
			}

		}

	}
}
