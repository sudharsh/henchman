package henchman

import (
	"fmt"
	"gopkg.in/yaml.v1"
	//"strings"
)

type TaskVars map[string]interface{}

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

func mergeMap(source *TaskVars, destination *TaskVars) {
	src := *source
	dst := *destination
	for variable, value := range src {
		dst[variable] = value
	}
}

/*
func PreparePlan(hostsFile string, sectionName string, planFile string) string {
	if hostsFile != nil {
		// read hostsFile
		hostsFileBuf, err := ioutil.ReadFile(hostsFile)
		if err != nil {
			fmt.Println("Error in reading hostsFile in PreparePlan")
			os.Exit(1)
		}

		// convert the yaml file into a map of string arrays
		m := make(map[interface{}][]string)
		err = yaml.Unmarshal(hostsFileBuf, &m)
	}

	// empty context struct
	ctxt := pongo2.Context{}

	// if there is a map of the host names, grab the specified section
	if m != nil && m[sectionName] != nil {
		nodeList := m[sectionName]
		ctxt = ctxt.Update(pongo2.Context{"nodes": nodeList})
	}

	tmpl, err := pongo2.FromFile(planFile)

	if err != nil {
		fmt.Println("Error in pongo2 read from file")
		os.Exit(1)
	}

	return tmpl.Execute(ctxt)
}
*/

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
	plan.parseTasks()
	return &plan, nil
}

func (plan *Plan) parseTasks() {
	for _, t := range plan.tasks {
		task := Task{}
		task.Name = t["name"]
		task.Action = t["action"]
		_, present := t["ignore_errors"]
		task.IgnoreErrors = bool(present)
		_, present = t["local"]
		task.LocalAction = bool(present)
		plan.Tasks = append(plan.Tasks, task)
	}
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
