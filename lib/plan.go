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
	Vars  TaskVars

	report map[string]string
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

func NewPlan(planBuf []byte, overrides TaskVars) (*Plan, error) {
	plan := Plan{}
	err := yaml.Unmarshal(planBuf, &plan)
	if err != nil {
		return nil, err
	}
	mergeMap(&overrides, &plan.Vars)
	plan.report = make(map[string]string)
	return &plan, nil
}

func (plan *Plan) PrintReport() {
	var counts = make(map[string]int)
	for _, status := range plan.report {
		_, present := counts[status]
		if !present {
			counts[status] = 1
		} else {
			counts[status]++
		}
	}
	fmt.Printf("Plan Report:\n")
	for k, v := range counts {
		fmt.Printf("%s - %d\n", k, v)
	}
}

func (plan *Plan) SaveStatus(task *Task, status string) {
	plan.report[task.Id] = status
}
