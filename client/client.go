// A NodNod client package.
package client

import (
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
	url       url.URL
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

	c.url = url.URL{Scheme: "ws", Host: c.Server, Path: "/"}

	dialer := websocket.Dialer{}

	conn, _, err := dialer.Dial(c.url.String(), nil)
	if err != nil {
		return &ConnectionErr{err}
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
	if c.connected == false {
		return nil, &NotConnectedErr{}
	}

	return c.stat(stats.STAT)
}

// Get stats of all NodNod servers/cluster.
func (c *Client) StatAll() (stats.Stats, error) {
	if c.connected == false {
		return nil, &NotConnectedErr{}
	}

	return c.stat(stats.STATALL)
}

// Check if client is connected.
func (c *Client) Connected() bool {
	return c.connected
}

func (c *Client) URL() string {
	return c.url.String()
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
				return nil, &StatFailedErr{resp.Error}
			}
			// All good
			return &resp, nil
		case err := <-c.err:
			return nil, err
		case <-timeout:
			return nil, &TimeoutErr{}
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
			c.err <- err
			break
		}
	}
}
