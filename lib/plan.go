package henchman

import (
	"fmt"
	"gopkg.in/yaml.v1"
)

type TaskVars map[string]string

type Plan struct {
	Hosts  []string
	Tasks  []Task
	Vars   TaskVars
	Report chan string
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
	plan.Report = make(chan string, len(plan.Tasks)*len(plan.Hosts))
	return &plan, nil
}

func (plan *Plan) PrintReport() {
	var report = make(map[string]int)
	for status := range plan.Report {
		_, present := report[status]
		if !present {
			report[status] = 1
		} else {
			report[status]++
		}
	}
	fmt.Printf("Plan Report:\n")
	for k, v := range report {
		fmt.Printf("%s - %d\n", k, v)
	}
}
