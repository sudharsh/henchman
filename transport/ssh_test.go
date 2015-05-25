package transport

import (
	"testing"
)

func TestValidPasswordAuth(t *testing.T) {
	c := make(TransportConfig)
	c["username"] = "user1"
	c["password"] = "password"
	c["hostname"] = "localhost"
	_, err := NewSSH(&c)
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestInvalidPasswordAuth(t *testing.T) {
	c := make(TransportConfig)
	c["username"] = "user1"
	c["hostname"] = "localhost"
	_, err := NewSSH(&c)
	if err == nil {
		t.Errorf("There should have been an error since password isn't present")
	}
}
