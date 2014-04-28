package henchman

import (
	"log"
	"fmt"
	"bytes"
	"code.google.com/p/go.crypto/ssh"
	"github.com/sudharsh/henchman/ansi"
)


type password string
func (p password) Password(pass string) (string, error) {
	return string(p), nil
}

type Machine struct {
	Username string
	Password string
	Hostname string
}

const (
	ECHO          = 53
	TTY_OP_ISPEED = 128
    TTY_OP_OSPEED = 129
)

func (machine *Machine) RunTask(task Task) {

	green := ansi.ColorCode("green")
    reset := ansi.ColorCode("reset")

	config := &ssh.ClientConfig{
		User: machine.Username,
	    Auth: []ssh.ClientAuth{
			ssh.ClientAuthPassword(password(machine.Password)),
		},
	}

	client, err := ssh.Dial("tcp", machine.Hostname + ":22", config)
	if err != nil {
		log.Fatalf("Failed to dial: " + err.Error())
	}

	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		log.Fatalf("Unable to create session: " + err.Error())
	}
	defer session.Close()

	modes := ssh.TerminalModes {
		ECHO: 0,
		TTY_OP_ISPEED: 14400,
		TTY_OP_OSPEED: 14400,
	}
	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		log.Fatalf("request for pseudo terminal failed: " + err.Error())
	}

	var b bytes.Buffer
	session.Stdout = &b
	log.Printf("**** Running task: %s\n", task["name"])
	log.Printf("---- Cmd: %s\n", task["action"])
	log.Printf("---- Host: %s\n", machine.Hostname)
	if err := session.Run(task["action"]); err != nil {
		panic("Failed to run: " + err.Error())
	}
	log.Printf("---- Output: \n")
	fmt.Print(green + b.String() + reset)
	log.Print("--------------------\n\n")
}
