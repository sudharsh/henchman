package henchman

import (
	"bytes"
	"code.google.com/p/go.crypto/ssh"
	"fmt"
	"github.com/sudharsh/henchman/ansi"
	"log"
	"text/template"
)

type Task struct {
	Name   string
	Action string
}

func prepareTemplate(data string, vars TaskVars) (string, error) {
	var buf bytes.Buffer
	tmpl, err := template.New("test").Parse(data)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(&buf, vars)
	return string(buf.Bytes()), err
}

func (task *Task) prepare(vars TaskVars) {
	var err error
	task.Name, err = prepareTemplate(task.Name, vars)
	if err != nil {
		panic(err)
	}
	task.Action, err = prepareTemplate(task.Action, vars)
	if err != nil {
		panic(err)
	}
}

func (task *Task) RunOn(machine *Machine, vars TaskVars, status chan string) {
	task.prepare(vars)
	green := ansi.ColorCode("green")
	red := ansi.ColorCode("red")
	reset := ansi.ColorCode("reset")

	client, err := ssh.Dial("tcp", machine.Hostname+":22", machine.SSHConfig)
	if err != nil {
		log.Fatalf("Failed to dial: " + err.Error())
	}

	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		log.Fatalf("Unable to create session: " + err.Error())
	}
	defer session.Close()

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
	log.Printf("**** Running task: %s\n", task.Name)
	log.Printf("---- Cmd: %s\n", task.Action)
	log.Printf("---- Host: %s\n", machine.Hostname)

	if err := session.Run(task.Action); err != nil {
		log.Printf("Failed to run: " + err.Error())
		fmt.Print(red + b.String() + reset)
		status <- "failure"
	} else {
		log.Printf("---- Success: \n")
		fmt.Print(green + b.String() + reset)
		status <- "success"
	}
	log.Print("--------------------\n\n")
}
