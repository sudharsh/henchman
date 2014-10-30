package henchman

import (
	"fmt"
	"gopkg.in/yaml.v1"
	"strings"
)

type TaskVars map[string]interface{}

// A plan is a collection of tasks.
// All the tasks are executed serially, although the same plan
// is run concurrently on multiple machines
type Plan struct {
	Hosts []string
	Tasks []Task
	Vars  *TaskVars
	Name  string

	report map[string]string
	tasks  []map[string]string `yaml:"tasks"`
}

func mergeMap(source *TaskVars, destination *TaskVars) {
	src := *source
	dst := *destination
	for variable, value := range src {
		dst[variable] = value
	}
}

// Returns a new plan with a collection of tasks. The planBuf should be a valid
// 'yaml' representation. Additionally, this function also takes in any variable
// overrides that takes precedence over the variables present in the plan.
func NewPlanFromYAML(planBuf []byte, overrides *TaskVars) (*Plan, error) {
	plan := Plan{}
	err := yaml.Unmarshal(planBuf, &plan)
	if plan.Vars == nil {
		_vars := make(TaskVars)
		plan.Vars = &_vars
	}
	if err != nil {
		return nil, err
	}
	if overrides != nil {
		mergeMap(overrides, plan.Vars)
		if hosts, present := (*overrides)["hosts"]; present {
			plan.Hosts = strings.Split(hosts.(string), ",")
		}

	}
	plan.report = make(map[string]string)
	plan.parseTasks()
	return &plan, nil
}

func (plan *Plan) parseTasks() {
	for _, t := range plan.tasks {
		task := Task{}
		task.Name = t["name"]
		task.Action = t["action"]
		_, present := t["ignore_errors"]
		task.IgnoreErrors = bool(present)
		_, present = t["local"]
		task.LocalAction = bool(present)
		plan.Tasks = append(plan.Tasks, task)
	}
}

// Prints the summary of the Plan execution across all the hosts
func (plan *Plan) PrintReport() {
	var counts = make(map[string]int)

	total := len(plan.Tasks) * len(plan.Hosts)
	attempted := len(plan.report)
	skipped := total - attempted

	counts["skipped"] = skipped
	for _, status := range plan.report {
		_, present := counts[status]
		if !present {
			counts[status] = 1
		} else {
			counts[status]++
		}
	}
	fmt.Println()
	fmt.Println("---")
	fmt.Printf("Plan Report: %s\n", plan.Name)
	fmt.Println()
	for k, v := range counts {
		fmt.Printf("%s (all hosts):\t%d\n", k, v)
	}
	fmt.Println()
	fmt.Printf("Tasks total (all hosts):\t%d\n",
		len(plan.Tasks)*len(plan.Hosts))
	fmt.Printf("Tasks attempted (all hosts):\t%d\n", len(plan.report))
}

// Mark a given task's status.
// NOTE: Skipped tasks are not tracked here.
func (plan *Plan) SaveStatus(task *Task, status string) {
	plan.report[task.Id] = status
}

func (plan *Plan) String() string {
	status := "Plan '%s':"
	for k, v := range plan.report {
		status = status + fmt.Sprintf(" %s - %s;", k, v)
	}
	return status
}
