package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/mohabusama/nodnod/client"
	"github.com/mohabusama/nodnod/stats"
)

type Node struct {
	name    string
	address string
	client  *client.Client
}

func newNode(address string) *Node {
	return &Node{
		address: address,
		client:  client.NewClient(address),
	}
}

func (n *Node) Connect() {
	if n.client.Connected() {
		log.Debug("Not connecting, Node already connected!")
		return
	}

	if err := n.client.Connect(); err == nil {
		log.Info("Established connection with node: ", n.address)
		// This should be able to load the name!
		n.Stat()
	} else {
		log.Warnf("Failed to connect to node: %s. Error: %s", n.address, err)
	}
}

func (n *Node) Stat() (*stats.NodeStats, error) {
	if st, err := n.client.Stat(); err != nil {
		return nil, err
	} else {
		// Return stats.NodeStat from stats.Stats
		for nodeName, nodeStats := range st {
			if n.name == "" {
				n.name = nodeName
				log.Info("Connected to node: ", n.String())
			}
			return &nodeStats, nil
		}
	}

	return &stats.NodeStats{}, nil
}

func (n *Node) Address() string {
	return n.address
}

func (n *Node) Name() string {
	if n.name == "" {
		return n.address
	}

	return n.name
}

func (n *Node) Connected() bool {
	return n.client.Connected()
}

func (n *Node) String() string {
	return fmt.Sprintf("[%s : %s]", n.Name(), n.address)
}

func (n *Node) Disconnect() {
	n.client.Disconnect()
}
