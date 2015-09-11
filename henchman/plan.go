package henchman

import (
	"fmt"
	//"github.com/flosch/pongo2"
	"gopkg.in/yaml.v1"
	"io/ioutil"
	//"strings"
)

type TaskVars map[string]interface{}

// A plan is a collection of tasks.
// All the tasks are executed serially, although the same plan
// is run concurrently on multiple machines
type Plan struct {
	Hosts []string
	Tasks []Task `yaml:"tasks"`
	Vars  *TaskVars
	Name  string

	report map[string]string
	//tasks  []map[string]string `yaml:"tasks"`
}

func mergeMap(source *TaskVars, destination *TaskVars) {
	src := *source
	dst := *destination
	for variable, value := range src {
		dst[variable] = value
	}
}

// Returns a new plan with a collection of tasks. The planBuf and hostsFileBuf should be a valid
// 'yaml' representation. Additionally, this function also takes in any variable
// overrides that takes precedence over the variables present in the plan.
func NewPlanFromYAML(planBuf []byte, hostsFileBuf []byte, overrides *TaskVars) (*Plan, error) {
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
		// if a hostsFile is specified, means checks the hosts override to get a list
		if hostsFileBuf != nil {
			if hosts, present := (*overrides)["hosts"]; present {
				hostsMap := make(map[interface{}][]string)
				err = yaml.Unmarshal(hostsFileBuf, &hostsMap)
				// should do something if there is an err, will get back to this
				// should I make this a plus equal?
				plan.Hosts = hostsMap[hosts]
			}
		}
	}
	plan.report = make(map[string]string)
	//fmt.Println(plan.Tasks)
	//plan.parseTasks()

	return &plan, nil
}

// change this to be a non-plan function
// come back to context variable stuff after getting include done
// look at diagram
/*
func (plan *Plan) PrepareTasks(tasks []Task, vars *TaskVars, machine Machine) ([]Task, error) {
	tmpl, err := pongo2.FromFile(planFile)
	if err != nil {
		return err
	}

	ctxt := pongo2.Context{"vars": plan.Vars, "machine": machine}

	out, err := tmpl.Execute(ctxt)
	if err != nil {
		return err
	}

	planBuf := []byte(out)

	newPlan := Plan{}
	err = yaml.Unmarshal(planBuf, &newPlan)
	if err != nil {
		return err
	}

	plan.Tasks = newPlan.Tasks

	return nil, nil
}
*/
// Updates the task list if there is a valid include param
// TODO: if there's a vars param too use the template in
//       the context of that vars,  currently it's using "global"
//       vars.
func UpdateTasks(tasks []Task, ndx int) ([]Task, error) {
	includeBuf, err := ioutil.ReadFile(tasks[ndx].Include)
	if err != nil {
		return nil, err
	}

	tmpPlan := Plan{}
	err = yaml.Unmarshal(includeBuf, &tmpPlan)
	if err != nil {
		return nil, err
	}

	//mergeMap(plan.Vars, tmpPlan.Vars)
	/*tmpPlan.Vars = plan.Vars
	err = tmpPlan.PrepareTasks(plan.Tasks[ndx].Include, machine)
	if err != nil {
		return err
	}*/

	newTasks := tmpPlan.Tasks
	tasks = append(tasks[:ndx+1], append(newTasks, tasks[ndx+1:]...)...)

	return tasks, nil
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
