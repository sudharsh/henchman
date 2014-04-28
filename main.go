package main

import (
	"flag"
	"log"
	"fmt"
	"os"
	"os/user"
	"code.google.com/p/gopass"
	"github.com/sudharsh/henchman/henchman"
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
	hostname  := flag.String("host", "localhost", "Host to run the module on")
	username  := flag.String("user", currentUsername(), "User to run as")
	password  := flag.String("password", "", "Path to the private key file")
	planFile := flag.String("plan", "", "Path to the plan")
	task := flag.String("task", "", "Task to run on the remote machine")

	flag.Parse()
	if *username == "" {
		os.Exit(1)
	}

	if *task == "" {
		log.Fatalf("Need a task")
		os.Exit(1)
	}

	if *password == "" {
		var err error
		if *password, err = gopass.GetPass("Password:"); err != nil {
			fmt.Println(err)
		}
	}

	plan, err := henchman.ParsePlan(*planFile)
	if err != nil {
		log.Fatalf("Couldn't read the plan: %s", err)
		os.Exit(1)
	}
	log.Printf("Fooo", plan.Tasks[0]["name"])
	var m = &henchman.Machine{*username, *password, *hostname, *task}
	m.RunTask()
}
