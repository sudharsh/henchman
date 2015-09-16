package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path"
	"strings"
	"sync"

	"code.google.com/p/gopass"

	"github.com/jlin21/henchman/henchman"
)

func currentUsername() *user.User {
	u, err := user.Current()
	if err != nil {
		log.Printf("Couldn't get current username: %s. Assuming root" + err.Error())
		u, err = user.Lookup("root")
		if err != nil {
			log.Print(err.Error())
		}
		return u
	}
	return u
}

func defaultKeyFile() string {
	u := currentUsername()
	return path.Join(u.HomeDir, ".ssh", "id_rsa")
}

// Split args from the cli that are of the form,
// "a=x b=y c=z" as a map of form { "a": "b", "b": "y", "c": "z" }
// These plan arguments override the variables that may be defined
// as part of the plan file.
func parseExtraArgs(args string) henchman.TaskVars {
	extraArgs := make(henchman.TaskVars)
	if args == "" {
		return extraArgs
	}
	for _, a := range strings.Split(args, " ") {
		kv := strings.Split(a, "=")
		extraArgs[kv[0]] = kv[1]
	}
	return extraArgs
}

// TODO: Modules
func validateModulesPath() (string, error) {
	_modulesDir := os.Getenv("HENCHMAN_MODULES_PATH")
	if _modulesDir == "" {
		cwd, _ := os.Getwd()
		_modulesDir = path.Join(cwd, "modules")
	}
	modulesDir := flag.String("modules", _modulesDir, "Path to the modules")
	_, err := os.Stat(*modulesDir)
	return *modulesDir, err
}

func localhost() *henchman.Machine {
	tc := make(henchman.TransportConfig)
	local, _ := henchman.NewLocal(&tc)
	localhost := henchman.Machine{}
	localhost.Hostname = "127.0.0.1"
	localhost.Transport = local
	return &localhost
}

func sshTransport(_tc *henchman.TransportConfig, hostname string) *henchman.SSHTransport {
	tc := *_tc
	tc["hostname"] = hostname
	ssht, err := henchman.NewSSH(&tc)
	if err != nil {
		log.Fatalf(err.Error())
	}
	return ssht
}

func main() {
	username := flag.String("user", currentUsername().Username, "User to run as")
	usePassword := flag.Bool("password", false, "Use password authentication")
	keyfile := flag.String("private-keyfile", defaultKeyFile(), "Path to the keyfile")
	extraArgs := flag.String("args", "", "Extra arguments for the plan")
	hostsFile := flag.String("i", "", "Specify hosts file name")

	modulesDir, err := validateModulesPath()
	if err != nil {
		log.Fatalf("Couldn't stat modules path '%s'\n", modulesDir)
	}
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [args] <plan>\n\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	planFile := flag.Arg(0)
	if planFile == "" {
		flag.Usage()
		os.Exit(1)
	}
	if *username == "" {
		fmt.Fprintf(os.Stderr, "Missing username")
		os.Exit(1)
	}
	// We support two SSH authentications methods for now
	// password and client key bases. Both are mutually exclusive and password takes
	// higher precedence
	tc := make(henchman.TransportConfig)
	tc["username"] = *username
	if *usePassword {
		var password string
		if password, err = gopass.GetPass("Password:"); err != nil {
			log.Fatalf("Couldn't get password: " + err.Error())
			os.Exit(1)
		}
		tc["password"] = password
	} else {
		tc["keyfile"] = *keyfile
	}

	planBuf, err := ioutil.ReadFile(planFile)
	if err != nil {
		log.Fatalf("Error reading plan - %s\n", planFile)
	}

	var hostsFileBuf []byte = nil
	if *hostsFile != "" {
		hostsFileBuf, err = ioutil.ReadFile(*hostsFile)
		if err != nil {
			log.Fatalf("Error reading hosts file - %s\n", *hostsFile)
		}
	}

	var plan *henchman.Plan
	parsedArgs := parseExtraArgs(*extraArgs)
	plan, err = henchman.NewPlanFromYAML(planBuf, hostsFileBuf, parsedArgs)
	if err != nil {
		log.Fatalf("Couldn't read the plan: %s", err)
		os.Exit(1)
	}

	local := localhost()
	// Execute the same plan concurrently across all the machines.
	// Note the tasks themselves in plan are executed sequentially.
	wg := new(sync.WaitGroup)
	for _, hostname := range plan.Hosts {
		machine := henchman.Machine{}
		machine.Hostname = hostname
		machine.Transport = sshTransport(&tc, hostname)

		// initializes a map for "register" values for each host
		registerMap := make(map[string]string)

		//renders all tasks in the plan file
		tasks, err := henchman.PrepareTasks(plan.Tasks, plan.Vars, machine)

		if err != nil {
			fmt.Println(err)
		}

		// for each host use the task list of the plan and run each task individually
		wg.Add(1)
		go func(machine *henchman.Machine) {
			defer wg.Done()

			// makes a temporary tasks to temper with
			for ndx := 0; ndx < len(tasks); ndx++ {
				var status *henchman.TaskStatus
				var err error
				task := tasks[ndx]

				// if there is a valid include field within task
				//    update task list
				// else
				//    do standard task run procedure
				whenVal, err := task.ProcessWhen(registerMap)
				if err != nil {
					log.Println("Error at When Eval at task: " + task.Name)
					log.Println("Error: " + err.Error())
				}
				if whenVal == true {
					if task.Include != "" {
						err = henchman.UpdateTasks(&tasks, task.Vars, ndx, *machine)
						if err != nil {
							log.Println("Error at Include Eval at task: " + task.Name)
							log.Println("Error: " + err.Error())
						}
					} else {
						if tasks[ndx].LocalAction {
							log.Printf("Local action detected\n")
							status, err = task.Run(local, registerMap)
						} else {
							status, err = task.Run(machine, registerMap)
						}
						plan.SaveStatus(&task, status.Status)
						if err != nil {
							log.Printf("Error when executing task: %s\n", err.Error())
						}
						if status.Status == "failure" {
							log.Printf("Task was unsuccessful: %s\n", task.Id)
							break
						}
					}
				}
			}
		}(&machine)
	}
	wg.Wait()
	plan.PrintReport()
}
