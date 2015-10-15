package client

import (
	"fmt"
)

type ConnectionErr struct {
	err error
}

type NotConnectedErr struct{}

type StatFailedErr struct {
	err string
}

type TimeoutErr struct{}

func (c *ConnectionErr) Error() string {
	return fmt.Sprintf("Failed to connect: %s", c.err.Error())
}

func (n *NotConnectedErr) Error() string {
	return "Client is not connected"
}

func (s *StatFailedErr) Error() string {
	return fmt.Sprintf("Stat failed: %s", s.err)
}

func (t *TimeoutErr) Error() string {
	return "Timeout error!"
}
