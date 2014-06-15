package henchman

import (
	"code.google.com/p/go.crypto/ssh"
)

type Machine struct {
	Hostname  string
	SSHConfig *ssh.ClientConfig
}


func Machines(hostnames []string, config *ssh.ClientConfig) []*Machine {
	var machines []*Machine
	for _, hostname := range hostnames {
		machine := Machine{hostname, config}
		machines = append(machines, &machine)
	}
	return machines
}

