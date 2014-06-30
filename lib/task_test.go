package henchman

import (
	"testing"
)

func TestPrepareTask(t *testing.T) {
	task := Task{"fake-uuid", "The {{ vars.variable1 }}", "The {{ vars.variable2 }}", false}
	vars := make(TaskVars)
	vars["variable1"] = "foo"
	vars["variable2"] = "bar"
	task.Prepare(&vars)
	if task.Name != "The foo" {
		t.Errorf("Template execution for Task.Name failed. Got - %s\n", task.Name)
	}
	if task.Action != "The bar" {
		t.Errorf("Template execution for Task.Action failed. Got - %s\n", task.Action)
	}
}
