package henchman

import (
	"bytes"
	"log"
	"os/exec"
	"strings"

	"code.google.com/p/go.crypto/ssh"
)

type Machine struct {
	Hostname  string
	SSHConfig *ssh.ClientConfig
}

func Machines(hostnames []string, config *ssh.ClientConfig) []*Machine {
	var machines []*Machine
	for _, hostname := range hostnames {
		machines = append(machines, &Machine{hostname, config})
	}
	return machines
}

// Exec this action on the machine
// TODO: Handle modules here
func (machine *Machine) Exec(action string) (*bytes.Buffer, error) {

	var b bytes.Buffer

	if machine.Hostname == "127.0.0.1" {
		log.Printf("Machines and action: %s\n", action)
		commands := strings.Split(action, " ")
		cmd := exec.Command(commands[0], commands[1:]...)
		cmd.Stdout = &b
		cmd.Stderr = &b
		err := cmd.Run()
		return &b, err
	}

	client, err := ssh.Dial("tcp", machine.Hostname+":22", machine.SSHConfig)
	if err != nil {
		log.Fatalf("Failed to dial: " + err.Error())
	}
	session, err := client.NewSession()
	if err != nil {
		log.Fatalf("Unable to create session: " + err.Error())
	}
	defer session.Close()
	defer client.Close()

	modes := ssh.TerminalModes{
		ECHO:          0,
		TTY_OP_ISPEED: 14400,
		TTY_OP_OSPEED: 14400,
	}
	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		log.Fatalf("request for pseudo terminal failed: " + err.Error())
	}
	session.Stdout = &b
	session.Stderr = &b
	return &b, session.Run(action)
}
