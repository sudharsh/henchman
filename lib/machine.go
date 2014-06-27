package henchman

import (
	"bytes"
	"log"

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

func (machine *Machine) RunTask(task *Task) string {
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

	var b bytes.Buffer
	session.Stdout = &b
	log.Printf("%s: %s '%s'\n", task.Id, machine.Hostname, task.Name)

	var taskStatus string = "success"
	if err := session.Run(task.Action); err != nil {
		if task.IgnoreErrors {
			taskStatus = "ignored"
		} else {
			taskStatus = "failure"
		}
	}
	escapeCode := statuses[taskStatus]
	var reset string = statuses["reset"]
	log.Printf("%s: %s [%s] - %s", task.Id, escapeCode, taskStatus, b.String()+reset)
	return taskStatus
}
