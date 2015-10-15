package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/mohabusama/nodnod/stats"
	"time"
)

// Collect stats form all connected NodNod servers.
func collect() ([]stats.NodeStats, error) {
	allStats := []stats.NodeStats{}
	chStat := make(chan *stats.NodeStats, len(globalNodes))

	for _, node := range globalNodes {

		go func(node *Node) {
			log.Debug("Stat for node: ", node.String())

			if st, err := node.Stat(); err == nil {
				log.Debug("Received stats from: ", node.String())
				chStat <- st
			} else {
				log.Debug("Received error from: ", node.String())
				chStat <- &stats.NodeStats{
					Name:  node.Name(),
					Error: err.Error(),
				}
			}
		}(node)
	}

	// All nodes stats should be available within 10 secs!
	timeout := time.After(10 * time.Second)
	statCount := len(globalNodes)
	timeoutErr := false

	for statCount > 0 && !timeoutErr {

		select {
		case nodeStat := <-chStat:
			log.Debug("Received stat channel: ", nodeStat)
			allStats = append(allStats, *nodeStat)
			statCount--
			log.Debug("Stats updated: ", allStats, statCount)
		case <-timeout:
			log.Warn("Get all stats timedout")
			timeoutErr = true
		}
	}

	return allStats, nil
}
