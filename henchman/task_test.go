package henchman

import (
	"testing"
)

// This function is no longer being used
/*
func TestPrepareTask(t *testing.T) {
	task := Task{"fake-uuid",
		"The {{ vars.variable1 }}",
		"{{ vars.variable2 }}:{{ machine.Hostname }}",
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
*/
func TestRun(t *testing.T) {
	task := Task{
		"fake-uuid",
		"The foo",
		"ls -al",
		"retVal",
		"",
		"",
		nil,
		false,
		false,
	}

	local, _ := NewLocal(nil)
	machine := Machine{"127.0.0.1", local}
	regMap := make(map[string]interface{})
	status, err := task.Run(&machine, regMap)
	if err != nil {
		t.Errorf("There shouldn't have been any error for this task")
	}
	if _, present := regMap["retVal"]; !present {
		t.Errorf("Register parameter should've had something stored")
	}
	if status.Status != "success" {
		t.Errorf("Task execution failed. Got %s\n", status)
	}
}

func TestProcessWhen(t *testing.T) {
	task := Task{
		"fake-uuid",
		"The foo",
		"echo hello",
		"",
		"",
		"first.Stdout == \"hello\"",
		nil,
		false,
		false,
	}

	output := RegisterOutput{
		Stdout: "hello",
		Stderr: "",
		Rc:     "success",
	}

	regMap := make(map[string]interface{})
	regMap["first"] = output

	whenVal, err := task.ProcessWhen(regMap)
	if err != nil {
		t.Errorf("There shouldn't have been any error using ProcessWhen")
	}
	if whenVal == false {
		t.Errorf("This ProcessWhen should always evaluate to true")
	}
}

func TestProcessWhenWithFilters(t *testing.T) {
	task := Task{
		"fake-uuid",
		"The foo",
		"echo hello",
		"",
		"",
		"first.Rc|success",
		nil,
		false,
		false,
	}

	output := RegisterOutput{
		Stdout: "hello",
		Stderr: "",
		Rc:     "success",
	}

	regMap := make(map[string]interface{})
	regMap["first"] = output

	whenVal, err := task.ProcessWhen(regMap)
	if err != nil {
		t.Errorf("There shouldn't have been any error using ProcessWhen")
	}
	if whenVal == false {
		t.Errorf("The ProcessWhen filter \"success\" should evaluate to true")
	}

	output.Rc = "failure"
	task.When = "first.Rc|failure"
	if whenVal == false {
		t.Errorf("The ProcessWhen filter \"failure\" should evaluate to true")
	}

	output.Rc = "ignored"
	task.When = "first.Rc|ignored"
	if whenVal == false {
		t.Errorf("The ProcessWhen filter \"ignored\" should evaluate to true")
	}

	output.Rc = "reset"
	task.When = "first.Rc|reset"
	if whenVal == false {
		t.Errorf("The ProcessWhen filter \"reset\" should evaluate to true")
	}
}
