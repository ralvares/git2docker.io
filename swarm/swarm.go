package swarm

import (
	"fmt"
	"github.com/ThomasRooney/gexpect"
	"github.com/cooltrick/cfg"
	"log"
	"os"
	"os/exec"
	"text/template"
)

type Swarmconf struct {
	Login        string
	Password     string
	Mail         string
	Ports        map[string]string
	Template     string
	Domain       string
	State        string
	Name         string
	Volumes      map[string]string
	Environments map[string]string
}

var (
	swarm, _        = exec.LookPath("swarm")
	userhome string = os.Getenv("HOME")
	username string = os.Getenv("USER")
)

func getSwarmBinaryPath() string {
	if len(swarm) != 0 {
		return swarm
	} else {
		return "error path -> Swarm"
	}
}

func (n *Swarmconf) GetInfos(name string) (string, string, string, string, map[string]string, string, string, string, map[string]string, map[string]string) {

	myconf := make(map[string]string)
	err := cfg.Load(os.TempDir()+"/"+username+"_"+name+"/swarm.conf", myconf)
	if err != nil {
		log.Fatal(err)
	}
	n.Name = name
	for k, v := range myconf {
		if k == "login" {
			n.Login = v
		}
		if k == "password" {
			n.Password = v
		}

		if k == "domain" {
			n.Domain = v
		}

		if k == "mail" {
			n.Mail = v
		}

	}

	return n.Login, n.Mail, n.Password, n.Domain, n.Ports, n.Template, n.State, n.Name, n.Volumes, n.Environments
}

func SwarmLogin(name string) bool {
	n := Swarmconf{}
	n.GetInfos(name)
	fmt.Printf("Login at swarm.. \n")
	child, err := gexpect.Spawn("swarm login  " + n.Mail)
	if err != nil {
		return false
	}

	child.Expect("Password:")
	child.SendLine(n.Password + "\n")
	succeeded := child.Expect("Login Succeeded")
	if succeeded != nil {
		child.Close()
		return false

	} else {
		child.Close()
		return true
	}

}

func SwarmrImageServer(name string) bool {
	n := Swarmconf{}
	n.GetInfos(name)
	fmt.Printf("Login at registry.giantswarm.io ... \n")

	imagesCmd := exec.Command("docker", "logout", "https://registry.giantswarm.io")
	out, err := imagesCmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(out))
		panic(err)
	}

	child, err := gexpect.Spawn("docker login -e " + n.Mail + " -u " + n.Login + " https://registry.giantswarm.io\n")
	if err != nil {
		return false
	}
	child.Expect("Password:")
	child.SendLine(n.Password + "\n")
	succeeded := child.Expect("Login Succeeded")
	if succeeded != nil {
		child.Close()
		return false

	} else {
		child.Close()
		return true
	}

}

func SwarmLS() {
	if len(swarm) > 0 {
		fmt.Println("Listing Dockers")
		imagesCmd := exec.Command(getSwarmBinaryPath(), "ls")
		out, err := imagesCmd.CombinedOutput()
		if err != nil {
			fmt.Println("Error - swarm ls")
			return
		} else {
			fmt.Println(string(out))
			return
		}

	}

}

func SwarmStop(name string) {
	n := Swarmconf{}
	n.GetInfos(name)

	if len(swarm) > 0 {
		fmt.Println("Stoping ->" + n.Name)
		imagesCmd := exec.Command(getSwarmBinaryPath(), "stop", n.Name)
		out, err := imagesCmd.CombinedOutput()
		if err != nil {
			fmt.Println("Error - swarm stop")
			return
		} else {
			fmt.Println(string(out))
			return
		}

	}

}

func SwarmStart(name string) {
	n := Swarmconf{}
	n.GetInfos(name)
	if len(swarm) > 0 {
		fmt.Println("Starting -> " + n.Name)
		imagesCmd := exec.Command(getSwarmBinaryPath(), "start", n.Name)
		out, err := imagesCmd.CombinedOutput()
		if err != nil {
			fmt.Println("Error - swarm start")
			return
		} else {
			fmt.Println(string(out))
			return
		}

	}

}

func SwarmImage(name string) bool {
	n := Swarmconf{}
	n.GetInfos(name)
	fmt.Println("Tag -> " + n.Name)
	RemoveTagCmd := exec.Command("docker", "rmi", "registry.giantswarm.io/"+n.Login+"/"+n.Name)
	RemoveTagCmd.Run()

	imagesCmd := exec.Command("docker", "tag", username+"/"+n.Name+":start", "registry.giantswarm.io/"+n.Login+"/"+n.Name)
	out, err := imagesCmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error - Tag Image")
		return false
	} else {
		fmt.Println(string(out))
		return true
	}
}

func PushImage(name string) bool {
	n := Swarmconf{}
	n.GetInfos(name)
	fmt.Println("Pushing Image -> " + n.Name)

	cmd := exec.Command("docker", "push", "registry.giantswarm.io/"+n.Login+"/"+n.Name)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println("ERROR - Push Image")
		return false
		panic(err)
	} else {
		return true
	}
}

func SwarmDelete(name string) {
	n := Swarmconf{}
	n.GetInfos(name)
	if len(swarm) > 0 {
		fmt.Println("Deleting -> " + n.Name)
		cmd := exec.Command(getSwarmBinaryPath(), "delete", "-y", n.Name)
		err := cmd.Run()
		if err != nil {
			fmt.Println("ERROR - Deleting Image")
			panic(err)
		}

	}

}

func SwarmUP(name string) {
	n := Swarmconf{}
	n.GetInfos(name)
	if len(swarm) > 0 {
		fmt.Println("Up -> " + n.Name)
		imagesCmd := exec.Command(getSwarmBinaryPath(), "up", os.Getenv("HOME")+"/"+name+"/swarm.json")
		imagesCmd.Stdout = os.Stdout
		imagesCmd.Stderr = os.Stderr
		err := imagesCmd.Run()
		if err != nil {
			fmt.Println("Error - swarm UP")
		}
	}

}

func SwarmGenerateJson(name string) bool {

	n := Swarmconf{}
	n.GetInfos(name)

	//tmpl, err := template.ParseFiles("/opt/git2docker/templates/" + templatename + "/swarm.tlp")
	tmpl, err := template.New("").Delims("[[[", "]]]").ParseFiles("/opt/git2docker/swarm.tlp")
	//tmpl, err := template.ParseGlob("*.tlp")
	if err != nil {
		panic(err)
	}

	swarmfile := os.Getenv("HOME") + "/" + name + "/swarm.json"
	filedocker, err := os.OpenFile(swarmfile, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	err = tmpl.ExecuteTemplate(filedocker, "swarm.tlp", n)
	if err != nil {
		return false
		panic(err)
	} else {
		return true
	}
}
