// Example of using raw Gorilla websocket dialer to stat NodNod server.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/websocket"
	"github.com/mohabusama/nodnod/stats"
	"net/url"
	"time"
)

const (
	MEGA = 1024 * 1024
)

var (
	flAddress = flag.String("server", "127.0.0.1:7070", "NodNod server to connect")

	flDuration = flag.Int("duration", 5, "Duration in seconds")
)

func main() {
	flag.Parse()

	log.SetLevel(log.InfoLevel)

	u := url.URL{Scheme: "ws", Host: *flAddress, Path: "/"}

	dialer := websocket.Dialer{}

	conn, _, err := dialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("Failed to establish connection with server: ", err)
	}

	defer conn.Close()

	now := time.Now()

	go func() {
		defer conn.Close()
		for {
			var mresp stats.MessageResponse

			if err := conn.ReadJSON(&mresp); err != nil {
				log.Error("Error reading json response: ", err)
				break
			} else {
				log.Info("Received response from node:", mresp.Host)
				log.Info("Duration:", time.Since(now))
				PrintStats(&mresp)
			}
		}

	}()

	for {
		var mreq stats.MessageRequest
		mreq.StatType = stats.ALL
		mreq.Type = stats.STATALL

		if err := conn.WriteJSON(&mreq); err != nil {
			log.Error("Failed to send req: ", err)
		} else {
			now = time.Now()
			log.Debug("Sent req: ", mreq)
		}

		time.Sleep(time.Duration(*flDuration) * time.Second)
	}
}

func PrintStats(mresp *stats.MessageResponse) {
	if prettyPrint, err := json.MarshalIndent(mresp, "", "    "); err == nil {
		fmt.Println(string(prettyPrint))
	} else {
		log.Error("Failed to unmarshal response: ", err)
	}
}
