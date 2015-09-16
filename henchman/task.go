package henchman

import (
	//"fmt"
	"log"
	"strconv"
	"strings"

	"code.google.com/p/go-uuid/uuid"
	"github.com/flosch/pongo2"

	"github.com/sudharsh/henchman/ansi"
)

var statuses = map[string]string{
	"reset":   ansi.ColorCode("reset"),
	"success": ansi.ColorCode("green"),
	"ignored": ansi.ColorCode("yellow"),
	"failure": ansi.ColorCode("red"),
}

type TaskStatus struct {
	Status  string
	Message string
}

type RegisterOutput struct {
	Stdout string
	Stderr string
	Rc     string
}

// Task is the unit of work in henchman.
type Task struct {
	Id           string
	Name         string
	Action       string
	Register     string
	Include      string
	When         string
	Vars         TaskVars
	IgnoreErrors bool `yaml:"ignore_errors"`
	LocalAction  bool `yaml:"local"`
}

/*
func prepareTemplate(data string, vars *TaskVars, machine *Machine) (string, error) {
	tmpl, err := pongo2.FromString(data)
	if err != nil {
		panic(err)
	}
	ctxt := pongo2.Context{"vars": vars, "machine": machine}
	return tmpl.Execute(ctxt)
}

// Renders the template parts in the task field.
// Also assigns a new UUID to the task uniquely identifying it.
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
*/
// Runs the task on the machine. The task might mutate `vars` so that other
// tasks down the `plan` can see any additions/updates.
func (task *Task) Run(machine *Machine, regMap map[string]interface{}) (*TaskStatus, error) {
	//task.prepare(vars, machine)
	task.Id = uuid.New()
	out, err := machine.Transport.Exec(task.Action)

	var taskStatus string = "success"
	var regOut RegisterOutput

	if err != nil {
		if task.IgnoreErrors {
			taskStatus = "ignored"
		} else {
			taskStatus = "failure"
		}
		regOut = RegisterOutput{
			Stderr: strings.TrimSpace(out.String()),
			Rc:     taskStatus,
		}
	} else {
		regOut = RegisterOutput{
			Stdout: strings.TrimSpace(out.String()),
			Rc:     taskStatus,
		}
	}

	if task.Register != "" {
		regMap[task.Register] = regOut
	}

	status := TaskStatus{taskStatus, out.String()}
	escapeCode := statuses[taskStatus]
	var reset string = statuses["reset"]
	log.Printf("%s: %s [%s] - %s", task.Id, escapeCode, status.Status, status.Message+reset)
	return &status, err
}

// Processes When param of a task.  Evaluates it using pongo2 and uses
// a registerMap (populated by other tasks register variables) as the
// context
func (task *Task) ProcessWhen(regMap map[string]interface{}) (bool, error) {
	if task.When == "" {
		return true, nil
	}

	tmpl, err := pongo2.FromString("{{ " + task.When + " }}")
	if err != nil {
		return false, err
	}

	// create context and execute
	/*
		ctxt := pongo2.Context{}
		for key, val := range regMap {
			ctxt = ctxt.Update(pongo2.Context{key: strings.TrimSpace(val)})
		}
	*/
	out, err := tmpl.Execute(regMap)

	if err != nil {
		return false, err
	}

	retVal, err := strconv.ParseBool(out)
	if err != nil {
		return false, err
	}

	return retVal, err
}
