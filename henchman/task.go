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
	Name         string
	Action       string
	IgnoreErrors bool `yaml:"ignore_errors"`
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

func (task *Task) Prepare(vars TaskVars) {
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
	task.Prepare(vars)
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

	var escapeCode string

	if err := session.Run(task.Action); err != nil {
		if task.IgnoreErrors {
			log.Printf("Ignoring this task's error")
			status <- "ignored"
			escapeCode = ansi.ColorCode("yellow")
		} else {
			log.Printf("Failed to run: " + err.Error())
			status <- "failure"
			escapeCode = ansi.ColorCode("red")
		}
	} else {
		log.Printf("---- Success: \n")
		escapeCode = ansi.ColorCode("green")
		status <- "success"
	}
	fmt.Print(escapeCode + b.String() + reset)
	log.Print("--------------------")
}
