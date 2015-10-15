package main

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"io/ioutil"
)

type Configuration struct {
	Nodes []string `json:"nodes"`

	// Alerts map[string]interface{} `json:"alerts"`

	Mode string `json:"mode"`
}

// Load configuration into globalConfig.
func loadConfig() error {
	data, err := ioutil.ReadFile(*flConfigPath)
	if err != nil {
		log.Error("Failed to read config file:", err)
		return err
	}

	err = json.Unmarshal(data, globalConfig)
	if err != nil {
		log.Error("Failed to load config data:", err)
		return err
	}

	return nil
}

// Load all other nodes into globalNodes.
func loadNodes() error {
	for _, n := range globalConfig.Nodes {
		if n == *flAddress {
			// Skip current node!
			continue
		}

		node := newNode(n)
		globalNodes = append(globalNodes, node)
	}

	return nil
}
