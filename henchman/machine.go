package henchman

import (
	"github.com/sudharsh/henchman/transport"
)

type Machine struct {
	Hostname  string
	Transport transport.TransportInterface
}
