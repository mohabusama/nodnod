// A NodNod client package.
package client

import (
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/mohabusama/nodnod/stats"
	"net/url"
	"sync"
	"time"
)

type Client struct {
	sync.Mutex

	conn      *websocket.Conn
	server    string
	connected bool
	response  chan stats.MessageResponse
	err       chan error
}

// Create a new client.
func NewClient(server string) *Client {
	return &Client{
		server:    server,
		connected: false,
	}
}

// Connect to NodNod server.
func (c *Client) Connect() error {
	if c.connected == true {
		return nil
	}

	u := url.URL{Scheme: "ws", Host: c.server, Path: "/"}

	dialer := websocket.Dialer{}

	conn, _, err := dialer.Dial(u.String(), nil)
	if err != nil {
		return err
	}

	c.conn = conn
	c.response = make(chan stats.MessageResponse, 1)
	c.err = make(chan error, 1)

	go c.listener()

	return nil
}

// Get all stats of NodNod cluster.
func (c *Client) Stat() (stats.AllStats, error) {

	mreq := stats.MessageRequest{
		Type:     stats.STATALL,
		StatType: stats.ALL,
	}

	c.Lock()
	if err := c.conn.WriteJSON(&mreq); err != nil {
		c.Unlock()
		return nil, err
	}

	c.Unlock()

	// Wait for response
	timeout := time.After(10 * time.Second)

	for {
		select {
		case resp := <-c.response:
			if resp.Error != "" {
				return nil, errors.New(fmt.Sprintf("Stat failed: %s", resp.Error))
			}
			// All good
			return resp.Nodes, nil
		case err := <-c.err:
			return nil, err
		case <-timeout:
			return nil, errors.New("Request timed out!")
		}
	}
}

// Check if client is connected.
func (c *Client) Connected() bool {
	return c.connected
}

// Disconnect client from NodNod server.
func (c *Client) Disconnect() {
	if c.connected == false {
		return
	}

	c.conn.Close()
	c.connected = false
}

// Establishes read loop to handle any response.
func (c *Client) listener() {
	defer c.Disconnect()

	for {
		var mresp stats.MessageResponse

		if err := c.conn.ReadJSON(&mresp); err == nil {
			c.response <- mresp
		} else {
			break
		}
	}
}
