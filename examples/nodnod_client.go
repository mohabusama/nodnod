// Example file using NodNod Client to stat a NodNod cluster.
package examples

import (
	"encoding/json"
	"flag"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/mohabusama/nodnod/client"
	"github.com/mohabusama/nodnod/stats"
	"time"
)

var (
	flServer = flag.String("server", "127.0.0.1:7070", "NodNod server address")
	flCount  = flag.Int("count", 10, "Number of requests. 0 means forever!")
)

func main() {

	flag.Parse()

	nodnodClient := client.NewClient(*flServer)

	if err := nodnodClient.Connect(); err != nil {
		log.Fatal("Failed to connect", err)
	}

	count := 0

	for {
		if allStats, err := nodnodClient.Stat(); err != nil {
			log.Error("Failed to get stats:", err)
		} else {
			prettyPrint(allStats)
		}

		count++
		if *flCount > 0 && count >= *flCount {
			log.Info("Reached maximum number of requests: ", count)
			nodnodClient.Disconnect()
			fmt.Println("Disconnected client.\nBye!")
			return
		}

		time.Sleep(5 * time.Second)
	}
}

func prettyPrint(allStats stats.AllStats) {
	if prettyPrint, err := json.MarshalIndent(allStats, "", "    "); err == nil {
		fmt.Println(string(prettyPrint))
	} else {
		log.Error("Failed to unmarshal response: ", err)
	}
}
