package henchman

import (
	"testing"
)

func TestPrepareTask(t *testing.T) {
	task := Task{"fake-uuid",
		"The {{ vars.variable1 }}",
		"{{ vars.variable2 }}:{{ machine.Hostname }}",
		false,
		false,
	}

	local, _ := NewLocal(nil)
	machine := Machine{"foobar", local}

	vars := make(TaskVars)
	vars["variable1"] = "foo"
	vars["variable2"] = "bar"
	task.prepare(&vars, &machine)

	if task.Name != "The foo" {
		t.Errorf("Template execution for Task.Name failed. Got - %s\n", task.Name)
	}
	if task.Action != "bar:foobar" {
		t.Errorf("Template execution for Task.Action failed. Got - %s\n", task.Action)
	}
}

func TestRun(t *testing.T) {
	task := Task{"fake-uuid",
		"The {{ vars.variable1 }}",
		"{{ vars.variable2 }}",
		false,
		false,
	}

	local, _ := NewLocal(nil)
	machine := Machine{"127.0.0.1", local}
	vars := make(TaskVars)

	vars["variable1"] = "foo"
	vars["variable2"] = "ls -al"

	status, err := task.Run(&machine, &vars)
	if err != nil {
		t.Errorf("There shouldn't have been any error for this task")
	}
	if status.Status != "success" {
		t.Errorf("Task execution failed. Got %s\n", status)
	}
}
