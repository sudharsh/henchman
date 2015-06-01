package henchman

import (
	"testing"
)

func TestValidLocalExec(t *testing.T) {
	c := make(TransportConfig)
	local, err := NewLocal(&c)
	if err != nil {
		t.Errorf(err.Error())
	}
	_, err = local.Exec("ls")
	if err != nil {
		t.Errorf(err.Error())
	}
}
