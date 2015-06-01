package henchman

import (
	"bytes"
	"log"
	"os/exec"
	"strings"
)

type LocalTransport struct{}

func (localTransport *LocalTransport) Initialize(config *TransportConfig) error {
	return nil
}

func (localTransport *LocalTransport) Exec(action string) (*bytes.Buffer, error) {
	var b bytes.Buffer
	log.Printf("action: %s\n", action)
	commands := strings.Split(action, " ")
	cmd := exec.Command(commands[0], commands[1:]...)
	cmd.Stdout = &b
	cmd.Stderr = &b
	err := cmd.Run()
	return &b, err
}

func NewLocal(config *TransportConfig) (*LocalTransport, error) {
	local := LocalTransport{}
	return &local, local.Initialize(config)
}
