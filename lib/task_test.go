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
	machine := Machine{"foobar", nil}

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
