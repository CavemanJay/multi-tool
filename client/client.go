package client

import (
	"fmt"
	"net/url"

	comms "github.com/JayCuevas/gogurt/communications"
	"github.com/gorilla/websocket"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("gogurt")

type Client struct {
	communicator *comms.Communicator
}

func NewClient() *Client {
	return &Client{
		communicator: comms.NewCommunicator(nil),
	}
}

func (c *Client) Connect(host string, port int) error {
	u := url.URL{
		Scheme: "ws",
		Host:   fmt.Sprintf("%s:%d", host, port),
		Path:   "/ws",
	}

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return err
	}

	defer conn.Close()

	c.communicator.HandleComms(conn)

	return nil
}
