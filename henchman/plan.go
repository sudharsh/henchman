package henchman

import (
	"log"
	"io/ioutil"
	"gopkg.in/yaml.v1"
)


type Plan struct {
	Hosts []string
	Tasks []map[string]string
}
		


func ParsePlan(hosts []string, config *string) (*Plan, error) {
	plan := Plan{}
	plan.Hosts = hosts
	data, err := ioutil.ReadFile(*config)
	log.Printf("%s", data)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(data, &plan)
	if err != nil {
		return nil, err
	}
	return &plan, nil
}