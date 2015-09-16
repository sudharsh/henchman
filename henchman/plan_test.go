package henchman

import "testing"

func TestParsePlanWithoutOverrides(t *testing.T) {
	plan_string := `---
name: "Sample plan"
hosts:
  - "127.0.0.1:22"
  - 192.168.1.2  
tasks:
  - name: Sample task that does nothing
    action: ls -al
  - name: Second task
    action: echo 'foo'
    ignore_errors: true
 `
	plan, err := NewPlanFromYAML([]byte(plan_string), nil, nil)
	if err != nil {
		panic(err)
	}
	if len(plan.Hosts) != 2 {
		t.Errorf("Number of hosts mismatch. Parsed %d hosts instead\n", len(plan.Hosts))
	}
	if len(plan.Tasks) != 2 {
		t.Errorf("Numnber of tasks mismatch. Parsed %d tasks instead\n", len(plan.Tasks))
	}
	if plan.Name != "Sample plan" {
		t.Errorf("Plan name mismath. Got %s\n", plan.Name)
	}
	second_task := plan.Tasks[1]
	if !second_task.IgnoreErrors {
		t.Errorf("The task '%s' had ignore_errors set to true. Got %t\n", second_task.Name, second_task.IgnoreErrors)
	}
}

func TestParsePlanWithHostOverrides(t *testing.T) {
	plan_string := `---
name: "Sample plan"
hosts:
  - "127.0.0.1:22"
  - 192.168.1.2  
tasks:
  - name: Sample task that does nothing
    action: ls -al
  - name: Second task
    action: echo 'foo'
    ignore_errors: true
 `

	hosts_string := `
group1:
  - 123.0.0.1
  - 123.0.0.2
  - 123.0.0.3
`
	tv := make(TaskVars)
	tv["hosts"] = "group1"

	plan, err := NewPlanFromYAML([]byte(plan_string), []byte(hosts_string), tv)
	if err != nil {
		panic(err)
	}
	if len(plan.Hosts) != 3 {
		t.Errorf("Number of hosts mismatch. Expecting 3. Parsed %d hosts instead\n", len(plan.Hosts))
	}
	if plan.Hosts[0] != "123.0.0.1" {
		t.Errorf("Hosts mismatch. Expecting 123.0.0.1. Parsed %s hosts instead\n", plan.Hosts[0])
	}
	if len(plan.Tasks) != 2 {
		t.Errorf("Numnber of tasks mismatch. Expected 2. Parsed %d tasks instead\n", len(plan.Tasks))
	}
	if plan.Name != "Sample plan" {
		t.Errorf("Plan name mismath. Expected Sample Plan. Got %s\n", plan.Name)
	}
}

func TestParsePlanWithOverrides(t *testing.T) {
	plan_string := `---
name: Sample plan
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

	plan, err := NewPlanFromYAML([]byte(plan_string), nil, tv)
	if err != nil {
		panic(err)
	}
	if len(plan.Hosts) != 1 {
		t.Errorf("Number of hosts mismatch. Parsed %d hosts instead\n", len(plan.Hosts))
	}
	if len(plan.Tasks) != 2 {
		t.Errorf("Numnber of tasks mismatch. Parsed %d tasks instead\n", len(plan.Tasks))
	}
	if plan.Name != "Sample plan" {
		t.Errorf("Plan name mismath. Got %s\n", plan.Name)
	}
	vars := plan.Vars
	if vars["service"] != "overridden_foo" {
		t.Error("Plan vars 'service' should have been 'overridden_foo'")
	}
	second_task := plan.Tasks[1]
	if !second_task.IgnoreErrors {
		t.Errorf("The task '%s' had ignore_errors set to false. Got %t\n", second_task.Name, second_task.IgnoreErrors)
	}
}

// Finish this test later.  We may change the current implementation of
// Prepare and UpdateTasks
/*
func TestParsePlanWithIncludes(t *testing.T) {
	plan_string := `---
name: "Sample plan"
hosts:
  - "127.0.0.1:22"
  - 192.168.1.2
tasks:
  - name: Sample task that does nothing
    action: ls -al

  - include: ../test/tasks.yaml

  - name: Second task
    action: echo 'foo'
    ignore_errors: true
 `
	hostbuf := []byte{}
	plan, err := NewPlanFromYAML([]byte(plan_string), hostbuf, nil)
	if err != nil {
		panic(err)
	}
	if len(plan.Hosts) != 2 {
		t.Errorf("Number of hosts mismatch. Expected 2. Parsed %d hosts instead\n", len(plan.Hosts))
	}
	if len(plan.Tasks) != 5 {
		t.Errorf("Numnber of tasks mismatch. Expected 5. Parsed %d tasks instead\n", len(plan.Tasks))
	}
	if plan.Name != "Sample plan" {
		t.Errorf("Plan name mismath. Got %s\n", plan.Name)
	}
	second_task := plan.Tasks[1]
}
*/
