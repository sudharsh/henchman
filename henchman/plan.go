package henchman

import (
	"fmt"
	"github.com/flosch/pongo2"
	"gopkg.in/yaml.v1"
	"io/ioutil"
	"strconv"
	//"strings"
)

type TaskVars map[string]interface{}

// A plan is a collection of tasks.
// All the tasks are executed serially, although the same plan
// is run concurrently on multiple machines
type Plan struct {
	Hosts []string
	Tasks []Task `yaml:"tasks"`
	Vars  TaskVars
	Name  string

	report map[string]string
}

// source values will override dest values override is true
// else dest values will not be overridden
func mergeMap(src TaskVars, dst TaskVars, override bool) {
	for variable, value := range src {
		if override == true {
			dst[variable] = value
		} else if _, present := dst[variable]; !present {
			dst[variable] = value
		}
	}
}

// Returns a new plan with a collection of tasks. The planBuf and hostsFileBuf should be a valid
// 'yaml' representation. Additionally, this function also takes in any variable
// overrides that takes precedence over the variables present in the plan.
func NewPlanFromYAML(planBuf []byte, hostsFileBuf []byte, overrides TaskVars) (*Plan, error) {
	plan := Plan{}
	err := yaml.Unmarshal(planBuf, &plan)

	if plan.Vars == nil {
		_vars := make(TaskVars)
		plan.Vars = _vars
	}
	if err != nil {
		return nil, err
	}

	if overrides != nil {
		mergeMap(overrides, plan.Vars, true)
		// if a hostsFile is specified, means checks the hosts override to get a list
		if hostsFileBuf != nil {
			if hosts, present := overrides["hosts"]; present {
				hostsMap := make(map[interface{}][]string)
				err = yaml.Unmarshal(hostsFileBuf, &hostsMap)
				// should do something if there is an err, will get back to this
				// should I make this a plus equal?
				plan.Hosts = hostsMap[hosts]
			}
		}
	}

	// populate Tasks Vars field with Plan Vars
	// or combine them only if task has a valid Include field
	// else task is a normal task and doesn't require Vars context
	// because tasks are rendered as a bundle and not individually
	// so vars context is only needed when calling PrepareTasks for Includes
	for ndx, task := range plan.Tasks {
		if task.Include != "" {
			if plan.Tasks[ndx].Vars != nil {
				mergeMap(plan.Vars, plan.Tasks[ndx].Vars, false)
			} else {
				plan.Tasks[ndx].Vars = plan.Vars
			}
		}
	}

	plan.report = make(map[string]string)

	return &plan, nil
}

// change this to be a non-plan function
// come back to context variable stuff after getting include done
// look at diagram
// renders Task list, uses vars and machine for context
func PrepareTasks(tasks []Task, vars TaskVars, machine Machine) ([]Task, error) {
	// changes Task array back to yaml form to be rendered
	var newTasks = []Task{}
	tasksBuf, err := yaml.Marshal(&tasks)
	if err != nil {
		return nil, err
	}

	// convert tasks to a pongo2 template
	tmpl, err := pongo2.FromString(string(tasksBuf))
	if err != nil {
		return nil, err
	}

	// add context and render
	ctxt := pongo2.Context{"vars": vars, "machine": machine}
	out, err := tmpl.Execute(ctxt)
	if err != nil {
		return nil, err
	}

	// change the newly rendered yaml format array back to a struct
	err = yaml.Unmarshal([]byte(out), &newTasks)
	if err != nil {
		return nil, err
	}

	return newTasks, nil
}

// Updates the task list if there is a valid include param
// for each include param it will update the vars context if
// it's provided.
func UpdateTasks(tasks []Task, vars TaskVars, ndx int, machine Machine) ([]Task, error) {
	includeBuf, err := ioutil.ReadFile(tasks[ndx].Include)
	if err != nil {
		return nil, err
	}

	// creates a new Plan object with only the tasks field filled in
	tmpPlan := Plan{}
	err = yaml.Unmarshal(includeBuf, &tmpPlan)
	if err != nil {
		return nil, err
	}

	// if there were no variables provided by the include task
	// populate it with variables in the previous context
	_vars := make(TaskVars)
	tmpPlan.Vars = _vars
	mergeMap(vars, tmpPlan.Vars, false)

	for ndx, task := range tmpPlan.Tasks {
		if task.Include != "" {
			if tmpPlan.Tasks[ndx].Vars != nil {
				mergeMap(tmpPlan.Vars, tmpPlan.Tasks[ndx].Vars, false)
			} else {
				tmpPlan.Tasks[ndx].Vars = tmpPlan.Vars
			}
		}
	}

	tmpPlan.Tasks, err = PrepareTasks(tmpPlan.Tasks, tmpPlan.Vars, machine)
	if err != nil {
		return nil, err
	}

	// insert the tasks in the tasks list
	tasks = append(tasks[:ndx+1], append(tmpPlan.Tasks, tasks[ndx+1:]...)...)

	return tasks, nil
}

// Evaluates the "When" parameter of a task using pongo2 templating
func CheckWhen(when string, regMap map[string]string) (bool, error) {
	if when == "" {
		return true, nil
	}

	tmpl, err := pongo2.FromString("{{ " + when + " }}")
	if err != nil {
		return false, err
	}

	// create context and execute
	ctxt := pongo2.Context{}
	for key, val := range regMap {
		ctxt = ctxt.Update(pongo2.Context{key: val[0 : len(val)-2]})
	}

	//ctxt = ctxt.Update(pongo2.Context{"first": "start"})
	fmt.Println(ctxt)
	out, err := tmpl.Execute(ctxt)
	fmt.Println("VALUE OF OUT IS: " + out)
	if err != nil {
		return false, err
	}

	retVal, err := strconv.ParseBool(out)
	if err != nil {
		return false, err
	}

	return retVal, nil
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
	status := fmt.Sprintf("Plan '%s' with %d tasks:", plan.Name, len(plan.Tasks))
	for k, v := range plan.report {
		status = status + fmt.Sprintf(" %s - %s;", k, v)
	}
	return status
}
