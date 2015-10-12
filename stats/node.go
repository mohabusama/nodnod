package stats

import (
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/websocket"
	"net/url"
	"sync"
)

type Node struct {
	conn   *websocket.Conn
	dialer websocket.Dialer
	mutex  *sync.Mutex

	Result  chan NodeStats
	Error   chan string
	Address string
	Status  bool
}

// Establish WS connection with a peer node.
func (this *Node) Connect() {
	if this.Status == true {
		log.Debug("Not connecting, Node already connected!")
		return
	}

	this.Result = make(chan NodeStats, 1)
	this.Error = make(chan string, 1)

	u := url.URL{Scheme: "ws", Host: this.Address, Path: "/"}

	this.dialer = websocket.Dialer{}

	conn, _, err := this.dialer.Dial(u.String(), nil)
	if err != nil {
		log.Warnf("Failed to connect to node: %s. Error: %s", this.Address, err)
		return
	}

	this.conn = conn
	defer this.Disconnect()

	this.Status = true
	this.mutex = new(sync.Mutex)

	log.Info("Established connection with node: ", this.Address)

	for {
		var mresp MessageResponse

		if err := this.conn.ReadJSON(&mresp); err == nil {
			if mresp.Error != "" {
				this.Error <- mresp.Error
			} else {
				if _, exists := mresp.Nodes[this.Address]; exists {
					this.Result <- mresp.Nodes[this.Address]
				} else {
					this.Error <- "Cannot find node address in response!"
				}
			}
		} else {
			log.Error("Error while reading: ", err)
			this.Error <- mresp.Error
			break
		}
	}
}

// Query node stat.
func (this *Node) Stat() error {
	if this.Status == false {
		// TODO: Define errors!
		return errors.New(
			fmt.Sprintf("Cannot stat node: %s. Node is not connected!", this.Address))
	}

	mReq := new(MessageRequest)
	mReq.Type = STAT
	mReq.StatType = ALL

	// Avoid concurrent calls to write!
	log.Debug("Acquiring lock!\n\n")
	this.mutex.Lock()
	log.Debug("Sending stat request via websocket \n\n")
	if err := this.conn.WriteJSON(&mReq); err != nil {
		this.mutex.Unlock()

		log.Error("Failed to send stat request to:", this.Address)

		return err
	} else {
		this.mutex.Unlock()
	}

	return nil
}

func (this *Node) Disconnect() {
	log.Info("Terminating node:", this.Address)

	defer func() {
		if r := recover(); r != nil {
			log.Debug("Recovered while terminating node:", r)
		}
	}()

	this.Status = false
	this.conn.Close()
}
