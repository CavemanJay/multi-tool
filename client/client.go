package client

import (
	"fmt"
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

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

	err = conn.WriteMessage(websocket.TextMessage, []byte("Test"))
	if err != nil {
		return err
	}

	for {
		_, bytes, err := conn.ReadMessage()
		if err != nil {
			panic(err)
		}

		log.Println(string(bytes))
	}

	return nil
}
