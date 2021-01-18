package communications

import (
	"fmt"
	"log"

	"github.com/JayCuevas/jays-server/sync"
	"github.com/gorilla/websocket"
)

type Communicator struct {
	conn         *websocket.Conn
	FileCreated  chan sync.File
	FilesDeleted chan []string
}

func NewCommunicator(socket *websocket.Conn) *Communicator {
	return &Communicator{
		conn:         socket,
		FileCreated:  make(chan sync.File),
		FilesDeleted: make(chan []string),
	}
}

func (c *Communicator) sendNewFileMessage(file *sync.File) error {
	event, err := fileCreatedEvent(file)
	if err != nil {
		return err
	}
	return c.conn.WriteJSON(event)
}

func (c *Communicator) handleNewFileMessage(e event) {

}

// TODO: Handle events not strings
func (c *Communicator) listenToClient() {
	for {
		var e event
		err := c.conn.ReadJSON(e)
		if err != nil {
			log.Println(err)
			return
		}

		switch e.id {
		case EventFileCreated:
			c.handleNewFileMessage(e)
		case EventFileDeleted:
		}
		// log.Println(string(p))
	}
}

func (c *Communicator) writeToClient() {
	for {
		select {
		case file := <-c.FileCreated:
			c.sendNewFileMessage(&file)
		case deletedFiles := <-c.FilesDeleted:
			fmt.Printf("Handling deleted files: %v", deletedFiles)
		}
	}
}

func (c *Communicator) HandleComms(conn *websocket.Conn) {
	c.conn = conn
	go c.writeToClient()

	c.listenToClient()
}
