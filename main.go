package main

import (
	"code.google.com/p/gopass"
	"code.google.com/p/go.crypto/ssh"
	"flag"
	"fmt"
	"github.com/sudharsh/henchman/henchman"
	"log"
	"os"
	"os/user"
)

func currentUsername() string {
	u, err := user.Current()
	if err != nil {
		log.Fatalf("Couldn't get current username")
		return ""
	}
	return u.Username
}

func main() {
	username := flag.String("user", currentUsername(), "User to run as")
	password := flag.String("password", "", "Password to the hosts")

	flag.Parse()
	planFile := flag.Arg(0)	
	if *username == "" {
		os.Exit(1)
	}

	if *password == "" {
		var err error
		if *password, err = gopass.GetPass("Password:"); err != nil {
			fmt.Println(err)
		}
	}

	plan, err := henchman.ParsePlan(planFile)
	if err != nil {
		log.Fatalf("Couldn't read the plan: %s", err)
		os.Exit(1)
	}

	sshAuth, err := henchman.SSHAuth(*password)
	if err != nil {
		log.Fatalf("SSH Auth prep failed: %s", err.Error())
	}
	config := &ssh.ClientConfig{
		User: *username,
	    Auth: sshAuth,
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
