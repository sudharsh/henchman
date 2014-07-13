package henchman

import "testing"

func TestMachines(t *testing.T) {
	hostnames := []string{
		"127.0.0.1",
		"192.168.33.11:8080",
	}
	machines := Machines(hostnames, nil)
	machine_one := machines[0]
	machine_two := machines[1]

	if machine_one.Hostname != "127.0.0.1" {
		t.Errorf("Hostname mismatched. Got %s hosts instead\n", machine_one.Hostname)
	}
	if machine_one.Port != 22 {
		t.Errorf("Port mismatched. Got %d hosts instead\n", machine_one.Port)
	}
	if machine_two.Hostname != "192.168.33.11" {
		t.Errorf("Hostname mismatched. Got %s hosts instead\n", machine_two.Hostname)
	}
	if machine_two.Port != 8080 {
		t.Errorf("Port mismatched. Got %d hosts instead\n", machine_two.Port)
	}
}
