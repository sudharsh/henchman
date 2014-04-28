package henchman

import (
	"log"
	"io/ioutil"
	"gopkg.in/yaml.v1"
)


type Plan struct {
	Tasks []map[string]string
}


func ParsePlan(config string) (*Plan, error) {
	plan := Plan{}
	data, err := ioutil.ReadFile(config)
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