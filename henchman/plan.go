package henchman

import (
	"gopkg.in/yaml.v1"
	"io/ioutil"
)

type TaskVars map[string]string

type Task struct {
	Name   string
	Action string
}

type Plan struct {
	Hosts []string
	Tasks []Task
	Vars  TaskVars
}

func mergeMap(source *TaskVars, dest *TaskVars) {
	s := *source
	d := *dest
	for variable, value := range s {
		d[variable] = value
	}
}

func ParsePlan(config string, overrides TaskVars) (*Plan, error) {
	plan := Plan{}

	data, err := ioutil.ReadFile(config)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(data, &plan)
	if err != nil {
		return nil, err
	}
	mergeMap(&overrides, &plan.Vars)
	return &plan, nil
}
