package henchman

import (
	"bytes"
	"gopkg.in/yaml.v1"
	"io/ioutil"
	"text/template"
)

type TaskVars map[string]string

type Task struct {
	Name   string
	Action string
}

type Plan struct {
	Hosts          []string
	Tasks          []Task
	Vars           TaskVars
	overriddenVars TaskVars
}

func mergeMap(source *TaskVars, dest *TaskVars) {
	s := *source
	d := *dest
	for variable, value := range s {
		d[variable] = value
	}
}

func prepareTemplate(config string, vars TaskVars) ([]byte, error) {
	var buf bytes.Buffer
	data, err := ioutil.ReadFile(config)
	tmpl, err := template.New("test").Parse(string(data[:]))
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(&buf, vars)
	return buf.Bytes(), err
}

func ParsePlan(config string, overrides TaskVars) (*Plan, error) {
	plan := Plan{}
	data, err := prepareTemplate(config, overrides)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(data, &plan)
	if err != nil {
		return nil, err
	}
	return &plan, nil
}
