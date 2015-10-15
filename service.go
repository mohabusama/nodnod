package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/websocket"
	"github.com/mohabusama/nodnod/stats"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  8192,
	WriteBufferSize: 8192,
}

type WebsocketHandler struct {
	http.Handler
}

// Handle websocket connections from clients.
func (this *WebsocketHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}

	conn, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Error("Failed to instantiate websocket:", err)
		return
	}

	defer conn.Close()

	log.Info("Accepted connection from client: ", conn.RemoteAddr())

	for {
		var mreq stats.MessageRequest
		mresp := stats.MessageResponse{
			Host:  *flName,
			Error: "",
			Nodes: map[string]stats.NodeStats{},
		}

		if err := conn.ReadJSON(&mreq); err == nil {

			switch mreq.Type {
			case stats.STAT:
				getCurrentNodeStat(&mresp)
				// Send response message

				if err = conn.WriteJSON(&mresp); err != nil {
					log.Error("Failed to respond with node stats:", err)
				}
			case stats.STATALL:
				getCurrentNodeStat(&mresp)
				getAllNodesStats(&mresp)

				log.Debug("STATALL nodes response:", mresp.Nodes)

				if err = conn.WriteJSON(&mresp); err != nil {
					log.Error("Failed to respond with all node stats:", err)
				}
			default:
				log.Warn("Received unknown request type:", mreq.Type)
			}
		} else {
			log.Debug("Terminated connection with client:", err)
			break
		}
	}
}

func getCurrentNodeStat(mresp *stats.MessageResponse) {
	if currentStat, err := stats.GetStats(*flName); err == nil {
		log.Debug("Current stat:", currentStat)

		mresp.Nodes[*flName] = currentStat
	} else {
		log.Error("Failed to get node stats:", err)
		mresp.Error = err.Error()
	}
}

func getAllNodesStats(mresp *stats.MessageResponse) {
	if len(globalNodes) == 0 {
		return
	}

	if allStats, err := collect(); err == nil {
		for _, nodeStat := range allStats {
			mresp.Nodes[nodeStat.Name] = nodeStat
		}
	} else {
		mresp.Error = err.Error()
	}
}
