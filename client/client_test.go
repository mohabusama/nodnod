package client

import (
	"testing"
)

var (
	server = "127.0.0.1:7070"
)

// Check connection failure.
func TestConnectionErr(t *testing.T) {
	cl := NewClient(server)

	if err := cl.Connect(); err == nil {
		t.Fatal("Expected error!")
	} else if _, ok := err.(*ConnectionErr); !ok {
		t.Errorf("Expected %s err, got %s instead.", &ConnectionErr{}, err)
	}
}

// Check connected status.
func TestInitialConnectedStatus(t *testing.T) {
	cl := NewClient(server)

	if cl.Connected() == true {
		t.Error("Expected connected status to be false, got: %s instead.", cl.Connected())
	}
}

// Check failed connected status.
func TestConnectedStatus(t *testing.T) {
	cl := NewClient(server)

	if err := cl.Connect(); err == nil {
		t.Fatal("Expected connect to fail!")
	}

	if cl.Connected() == true {
		t.Error("Expected connected status to be false, got: %s instead.", cl.Connected())
	}
}

// Attempt calling Stat without connect.
func TestStatNotConnectedErr(t *testing.T) {
	cl := NewClient(server)

	_, err := cl.Stat()
	if err == nil {
		t.Fatal("Expected error!")
	} else {
		if _, ok := err.(*NotConnectedErr); !ok {
			t.Errorf("Expected %s error, got %s", &NotConnectedErr{}, err)
		}
	}
}

// Attempt calling StatAll without connect.
func TestStatAllNotConnectedErr(t *testing.T) {
	cl := NewClient(server)

	_, err := cl.StatAll()
	if err == nil {
		t.Fatal("Expected error!")
	} else {
		if _, ok := err.(*NotConnectedErr); !ok {
			t.Errorf("Expected %s error, got %s", &NotConnectedErr{}, err)
		}
	}
}

// Test URL.
func TestURL(t *testing.T) {
	cl := NewClient(server)
	expected := "ws://127.0.0.1:7070/"

	cl.Connect()

	if cl.URL() != expected {
		t.Errorf("Expected url to be: %s, got %s instead.", expected, cl.URL())
	}
}

// Disconnect non-connected client
func TestDisconnectNonConnected(t *testing.T) {
	cl := NewClient(server)

	cl.Connect()

	cl.Disconnect()
	cl.Disconnect()
}
