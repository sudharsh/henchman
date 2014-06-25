package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path"
	"strings"
	"sync"

	"code.google.com/p/go.crypto/ssh"
	"code.google.com/p/gopass"

	"github.com/sudharsh/henchman/lib"
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

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [args] <plan>\n\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	planFile := flag.Arg(0)

	if planFile == "" {
		flag.Usage()
		os.Exit(1)
	}

	if *username == "" {
		fmt.Fprintf(os.Stderr, "Missing username")
		os.Exit(1)
	}

	var sshAuth ssh.AuthMethod
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
		Auth: []ssh.AuthMethod{sshAuth},
	}

	planBuf, err := ioutil.ReadFile(planFile)
	if err != nil {
		log.Fatalf("Error reading plan - %s\n", planFile)
		os.Exit(1)
	}

	plan, err := henchman.NewPlan(planBuf, parseExtraArgs(*extraArgs))
	if err != nil {
		log.Fatalf("Couldn't read the plan: %s", err)
		os.Exit(1)
	}

	wg := new(sync.WaitGroup)
	machines := henchman.Machines(plan.Hosts, config)
	for _, _machine := range machines {
		machine := _machine
		wg.Add(1)
		go func() {
			defer wg.Done()
			for _, task := range plan.Tasks {
				task.RunOn(machine, plan.Vars, plan.Report)
			}
		}()
	}
	wg.Wait()
	close(plan.Report)
	plan.PrintReport()
}
