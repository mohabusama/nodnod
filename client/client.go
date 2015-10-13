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

	Server string

	conn      *websocket.Conn
	connected bool
	response  chan stats.MessageResponse
	err       chan error
}

// Create a new client.
func NewClient(server string) *Client {
	return &Client{
		Server: server,
	}
}

// Connect to NodNod server.
func (c *Client) Connect() error {
	if c.connected == true {
		return nil
	}

	u := url.URL{Scheme: "ws", Host: c.Server, Path: "/"}

	dialer := websocket.Dialer{}

	conn, _, err := dialer.Dial(u.String(), nil)
	if err != nil {
		return err
	}

	c.conn = conn
	c.connected = true
	c.response = make(chan stats.MessageResponse, 1)
	c.err = make(chan error, 1)

	go c.read()

	return nil
}

// Get stats of *only* connected NodNod server.
func (c *Client) Stat() (stats.Stats, error) {
	return c.stat(stats.STAT)
}

// Get stats of all NodNod servers/cluster.
func (c *Client) StatAll() (stats.Stats, error) {
	return c.stat(stats.STATALL)
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

func (c *Client) stat(reqType int) (stats.Stats, error) {
	mreq := stats.MessageRequest{
		Type:     reqType,
		StatType: stats.ALL,
	}

	if mresp, err := c.write(mreq); err == nil {
		return mresp.Nodes, nil
	} else {
		return nil, err
	}
}

func (c *Client) write(mreq stats.MessageRequest) (*stats.MessageResponse, error) {
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
			return &resp, nil
		case err := <-c.err:
			return nil, err
		case <-timeout:
			return nil, errors.New("Request timed out!")
		}
	}
}

// Establishes read loop to handle any response.
func (c *Client) read() {
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
