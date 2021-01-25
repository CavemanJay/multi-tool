package client

import (
	"fmt"
	"net/url"

	comms "github.com/CavemanJay/gogurt/communications"
	"github.com/CavemanJay/gogurt/sync"
	"github.com/gorilla/websocket"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("gogurt")

type Client struct {
	communicator *comms.Communicator
	root         string
}

func NewClient(root string) *Client {
	c := comms.NewCommunicator(nil, "")
	c.FileReceived = onFileReceived
	return &Client{
		communicator: c,
		root:         root,
	}
}

func onFileReceived(file sync.FileWithData) {
	log.Debugf("Received file: %v", file)
}

func (c *Client) Connect(host string, port int) error {
	u := url.URL{
		Scheme: "ws",
		Host:   fmt.Sprintf("%s:%d", host, port),
		Path:   "/ws",
	}

	log.Debugf("Attempting to connect to %s", u.String())
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return err
	}

	defer conn.Close()

	c.communicator.HandleComms(conn)

	return nil
}
