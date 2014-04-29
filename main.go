package main

import (
	"code.google.com/p/go.crypto/ssh"
	"code.google.com/p/gopass"
	"flag"
	"github.com/sudharsh/henchman/henchman"
	"log"
	"os"
	"os/user"
	"path"
)

func currentUsername() *user.User {
	u, err := user.Current()
	if err != nil {
		panic("Couldn't get current username: " + err.Error())
	}
	return u
}

func defaultKeyFile() string {
	u := currentUsername()
	return path.Join(u.HomeDir, ".ssh", "id_rsa")
}

func main() {
	username := flag.String("user", currentUsername().Username, "User to run as")
	usePassword := flag.Bool("password", false, "Use password authentication")
	keyfile := flag.String("private-keyfile", defaultKeyFile(), "Path to the keyfile")

	flag.Parse()
	planFile := flag.Arg(0)
	if *username == "" {
		os.Exit(1)
	}

	var sshAuth ssh.ClientAuth
	var err error
	if *usePassword {
		var password string
		if password, err = gopass.GetPass("Password:"); err != nil {
			log.Fatalf("Couldn't get password: " + err.Error())
			os.Exit(1)
		}
		sshAuth, err = henchman.PasswordAuth(password)
	} else {
		sshAuth, err = henchman.ClientKeyAuth(*keyfile)
	}
	if err != nil {
		log.Fatalf("SSH Auth prep failed: " + err.Error())
	}
	config := &ssh.ClientConfig{
		User: *username,
		Auth: []ssh.ClientAuth{sshAuth},
	}

	plan, err := henchman.ParsePlan(planFile)
	if err != nil {
		log.Fatalf("Couldn't read the plan: %s", err)
		os.Exit(1)
	}

	sem := make(chan int, 100)
	for _, hostname := range plan.Hosts {
		go func() {
			machine := henchman.Machine{hostname, config}
			for _, task := range plan.Tasks {
				machine.RunTask(&task)
			}
			sem <- 1
		}()
		<-sem
	}

}
