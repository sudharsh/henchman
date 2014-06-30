package henchman

import (
	"fmt"
	"gopkg.in/yaml.v1"
)

// TaskVars hold any variables that is overridden through the cli
type TaskVars map[string]string

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

const (
	ECHO          = 53
	TTY_OP_ISPEED = 128
	TTY_OP_OSPEED = 129
)

func mergeMap(source *TaskVars, destination *TaskVars) {
	src := *source
	dst := *destination
	for variable, value := range src {
		dst[variable] = value
	}
}

func NewPlan(planBuf []byte, overrides *TaskVars) (*Plan, error) {
	plan := Plan{}
	err := yaml.Unmarshal(planBuf, &plan)
	if err != nil {
		return nil, err
	}
	if overrides != nil {
		mergeMap(overrides, plan.Vars)
	}
	plan.report = make(map[string]string)
	return &plan, nil
}

func (plan *Plan) parseTasks() {
	for _, t := range plan.tasks {
		task := Task{}
		task.Name = t["name"]
		task.Action = t["action"]
		_, present := t["ignore_errors"]
		task.IgnoreErrors = bool(present)
		plan.Tasks = append(plan.Tasks, task)
	}
}

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

func (plan *Plan) SaveStatus(task *Task, status string) {
	plan.report[task.Id] = status
}
