package main

import (
	log "github.com/Sirupsen/logrus"
	"time"
)

func discover() {
	log.Info("Starting discovery service")

	// Loop to make sure we keep all nodes connected!
	for {
		for _, node := range globalNodes {
			if node.Status == true {
				continue
			}

			go node.Connect()
		}

		time.Sleep(5 * time.Second)
	}
}
