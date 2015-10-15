package client

import (
	"testing"
)

// Attempt calling Stat/StatAll without connect.
func TestStatNotConnectedErr(t *testing.T) {
	cl := NewClient("127.0.0.1:7070")

	_, err := cl.Stat()
	if err == nil {
		t.Fatal("Expected error!")
	} else {
		if _, ok := err.(*NotConnectedErr); !ok {
			t.Errorf("Expected %s error, got %s", &NotConnectedErr{}, err)
		}
	}
}
