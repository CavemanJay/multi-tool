package client

import (
	"fmt"
	"net/url"

	"github.com/gorilla/websocket"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("gogurt")

type Client struct {
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

	// err = conn.WriteMessage(websocket.TextMessage, []byte("Test"))
	// if err != nil {
	// 	return err
	// }

	for {
		_, bytes, err := conn.ReadMessage()
		if err != nil {
			log.Error(err)
		}

		log.Debugf("Received message from server: %s", bytes)
	}

	return nil
}
