package main

import (
	"code.google.com/p/go.crypto/ssh"
	"code.google.com/p/gopass"
	"flag"
	"github.com/sudharsh/henchman/henchman"
	"log"
	"strings"

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

func parseExtraArgs(args string) map[string]string {
	extraArgs := make(map[string]string)
	if args == "" {
		return extraArgs
	}
	for _, a := range strings.Split(args, " ") {
		kv := strings.Split(a, "=")
		extraArgs[kv[0]] = kv[1]
	}
	return extraArgs
}

func main() {
	username := flag.String("user", currentUsername().Username, "User to run as")
	usePassword := flag.Bool("password", false, "Use password authentication")
	keyfile := flag.String("private-keyfile", defaultKeyFile(), "Path to the keyfile")
	extraArgs := flag.String("args", "", "Extra arguments for the plan")
	flag.Parse()

	args := parseExtraArgs(*extraArgs)
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
	for _, task := range plan.Tasks {
		for _, hostname := range plan.Hosts {
			go func() {
				machine := henchman.Machine{hostname, config}
				machine.RunTask(&task)
				sem <- 1
			}()
			<-sem
		}
	}
}
