package henchman

import (
	"io/ioutil"
	"gopkg.in/yaml.v1"
)


type Task map[string]string

type Plan struct {
	Hosts []string
	Tasks []Task
}
		
func ParsePlan(config *string) (*Plan, error) {
	plan := Plan{}

	data, err := ioutil.ReadFile(*config)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(data, &plan)
	if err != nil {
		return nil, err
	}
	return &plan, nil
}

