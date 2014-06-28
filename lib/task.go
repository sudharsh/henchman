package henchman

import (
	"bytes"
	"text/template"

	"code.google.com/p/go-uuid/uuid"

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
	IgnoreErrors bool
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
