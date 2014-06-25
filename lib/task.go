package henchman

import (
	"bytes"
	"log"
	"text/template"

	"code.google.com/p/go.crypto/ssh"
	"code.google.com/p/go-uuid/uuid"
	
	"github.com/sudharsh/henchman/ansi"
)

type Task struct {
	Id string
	
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
	task.Id = uuid.New()
	task.Name, err = prepareTemplate(task.Name, vars)
	if err != nil {
		panic(err)
	}
	task.Action, err = prepareTemplate(task.Action, vars)
	if err != nil {
		panic(err)
	}
}

func (task *Task) RunOn(machine *Machine, vars TaskVars) (string) {
	task.Prepare(vars)
	reset := ansi.ColorCode("reset")

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

	var escapeCode string
	var taskStatus string
	if err := session.Run(task.Action); err != nil {
		if task.IgnoreErrors {
			taskStatus = "ignored"
			escapeCode = ansi.ColorCode("yellow")
		} else {
			taskStatus = "failure"
			escapeCode = ansi.ColorCode("red")
		}
	} else {
		escapeCode = ansi.ColorCode("green")
		taskStatus = "success"
	}
	log.Printf("%s: %s [%s] - %s", task.Id, escapeCode, taskStatus, b.String() + reset)
	return taskStatus
}
