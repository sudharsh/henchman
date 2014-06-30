package henchman

import (
	"bytes"
	"log"

	"code.google.com/p/go-uuid/uuid"
	"code.google.com/p/go.crypto/ssh"
	"github.com/flosch/pongo2"

	"github.com/sudharsh/henchman/ansi"
)

var statuses = map[string]string{
	"reset":   ansi.ColorCode("reset"),
	"success": ansi.ColorCode("green"),
	"ignored": ansi.ColorCode("yellow"),
	"failure": ansi.ColorCode("red"),
}

// Task is the unit of work in henchman.
type Task struct {
	Id string

	Name         string
	Action       string
	IgnoreErrors bool `yaml:"ignore_errors"`
}

func prepareTemplate(data string, vars *TaskVars, machine *Machine) (string, error) {
	tmpl, err := pongo2.FromString(data)
	if err != nil {
		panic(err)
	}
	return tmpl.Execute(&pongo2.Context{"vars": vars, "machine": machine})
}

func (task *Task) prepare(vars *TaskVars, machine *Machine) {
	var err error
	task.Id = uuid.New()
	task.Name, err = prepareTemplate(task.Name, vars, machine)
	if err != nil {
		panic(err)
	}
	task.Action, err = prepareTemplate(task.Action, vars, machine)
	if err != nil {
		panic(err)
	}
}

func (task *Task) Run(machine *Machine, vars *TaskVars) string {
	task.prepare(vars, machine)
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
