package main

import (
	"flag"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/websocket"
	"github.com/mohabusama/nodnod/stats"
	"net/url"
	"os"
	"sync"
	"time"
)

var (
	flConcurrent = flag.Int("concurrent", 10, "Number of concurrent connections")

	flNumber = flag.Int("number", 100, "Total number of requests")

	// flRate = flag.Int("rate", 0, "Number of requests per seconds")

	flAddress = flag.String("server", "127.0.0.1:7070", "NodNod server to connect")

	flNodes = flag.Int("nodes", 1, "Number of expected nodes to stat. Used in validating the responses.")
)

type Response struct {
	Success  bool
	Duration time.Duration
	Valid    bool
}

func main() {

	flag.Parse()

	if *flNumber < *flConcurrent {
		log.Fatalf(
			"Total number of requests %d should be greater than number of concurrent connections %s",
			*flNumber, *flConcurrent)
		os.Exit(1)
	}

	flag.Usage = func() {
		flag.PrintDefaults()
		return
	}

	log.SetLevel(log.InfoLevel)

	u := url.URL{Scheme: "ws", Host: *flAddress, Path: "/"}

	chResponses := make(chan *Response, *flNumber)
	chRequests := make(chan stats.MessageRequest, *flNumber)

	var wg sync.WaitGroup
	wg.Add(*flNumber)

	for i := 0; i < *flConcurrent; i++ {
		go nodnodClient(&u, &wg, chRequests, chResponses, i)
	}

	log.Info("Launched all goroutines!")

	var mreq stats.MessageRequest
	mreq.StatType = stats.ALL
	mreq.Type = stats.STATALL

	start := time.Now()
	for i := 0; i < *flNumber; i++ {
		chRequests <- stats.MessageRequest{
			StatType: stats.ALL,
			Type:     stats.STATALL,
		}
	}

	log.Info("Launched all requests!")

	close(chRequests)

	log.Info("In progress ...")

	wg.Wait()

	printSummary(chResponses, time.Since(start))
}

func nodnodClient(url *url.URL, wg *sync.WaitGroup, chRequests chan stats.MessageRequest,
	chResponse chan *Response, id int) {

	dialer := websocket.Dialer{}

	conn, _, err := dialer.Dial(url.String(), nil)
	if err != nil {
		log.Fatal("Failed to establish connection with server: ", err)
	}

	defer conn.Close()
	start := time.Now()
	connected := true
	var clientWg sync.WaitGroup

	go func() {
		for connected {
			var mresp stats.MessageResponse

			if err := conn.ReadJSON(&mresp); err != nil {
				if connected {
					log.Error("Error reading json response: ", err)
					chResponse <- &Response{
						Success: false,
					}
					clientWg.Done()
					wg.Done()
				}
			} else {
				log.Debug("Received response Host:", mresp.Host)

				duration := time.Since(start)
				chResponse <- &Response{
					Success:  true,
					Duration: duration,
					Valid:    validateNodes(&mresp),
				}

				clientWg.Done()
				wg.Done()
			}
		}

	}()

	for mreq := range chRequests {
		clientWg.Add(1)

		if err := conn.WriteJSON(&mreq); err != nil {
			log.Error("Failed to send req: ", err)
			clientWg.Done()
		} else {
			start = time.Now()
			log.Debug("Sent req: ", mreq)
		}
	}

	clientWg.Wait()
	connected = false
}

func validateNodes(mresp *stats.MessageResponse) bool {
	return len(mresp.Nodes) == *flNodes
}

func printSummary(chResponses chan *Response, totalDuration time.Duration) {
	totalResponses := 0
	totalFailures := 0
	inValidResponses := 0

	waiting := true

	for waiting {
		select {
		case resp := <-chResponses:
			if resp.Success == false {
				// Failed response
				totalFailures++
			} else {
				if resp.Valid == false {
					// Invalid response
					inValidResponses++
				}
				totalResponses++
			}
		default:
			waiting = false
		}
	}

	log.Info("====SUMMARY====")

	log.Info("Total Duration:", totalDuration)
	log.Info("Total number of requests:", *flNumber)
	log.Info("Total number of responses:", totalResponses)
	log.Info("Total number of missing responses:", *flNumber-totalResponses)
	log.Info("Total number of failed responses:", totalFailures)
	log.Info("Total number of invalid responses:", inValidResponses)

	log.Info("========")
}
