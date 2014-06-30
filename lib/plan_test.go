package henchman

import "testing"

func TestParsePlanWithoutOverrides(t *testing.T) {
	plan_string := `---
hosts:
  - 127.0.0.1
tasks:
  - name: Sample task that does nothing
    action: ls -al
  - name: Second task
    action: echo 'foo'
    ignore_errors: true
 `
	plan, err := NewPlan([]byte(plan_string), nil)
	if err != nil {
		panic(err)
	}
	if len(plan.Hosts) != 1 {
		t.Errorf("Number of hosts mismatch. Parsed %d hosts instead\n", len(plan.Hosts))
	}
	if len(plan.Tasks) != 2 {
		t.Errorf("Numnber of tasks mismatch. Parsed %d tasks instead\n", len(plan.Tasks))
	}
	second_task := plan.Tasks[1]
	if !second_task.IgnoreErrors {
		t.Errorf("The task '%s' had ignore_errors set to true. Got %t\n", second_task.Name, second_task.IgnoreErrors)
	}
}

func TestParsePlanWithOverrides(t *testing.T) {
	plan_string := `---
vars:
  service: foo
hosts:
  - 127.0.0.1
tasks:
  - name: Sample task that does nothing
    action: ls -al
  - name: Second task
    action: echo 'foo'
    ignore_errors: true
 `
	tv := make(TaskVars)
	tv["service"] = "overridden_foo"

	plan, err := NewPlan([]byte(plan_string), &tv)
	if err != nil {
		panic(err)
	}
	if len(plan.Hosts) != 1 {
		t.Errorf("Number of hosts mismatch. Parsed %d hosts instead\n", len(plan.Hosts))
	}
	if len(plan.Tasks) != 2 {
		t.Errorf("Numnber of tasks mismatch. Parsed %d tasks instead\n", len(plan.Tasks))
	}
	vars := *plan.Vars
	if vars["service"] != "overridden_foo" {
		t.Error("Plan vars 'service' should have been 'overridden_foo'")
	}
	second_task := plan.Tasks[1]
	if !second_task.IgnoreErrors {
		t.Errorf("The task '%s' had ignore_errors set to false. Got %t\n", second_task.Name, second_task.IgnoreErrors)
	}
}
