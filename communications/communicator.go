package communications

import (
	"encoding/json"
	"strings"

	"github.com/JayCuevas/gogurt/sync"
	"github.com/gorilla/websocket"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("gogurt")

type Communicator struct {
	conn         *websocket.Conn
	syncFolder   string
	FileCreated  chan sync.FileWithData
	FilesDeleted chan []string
	FileReceived func(file sync.FileWithData)
}

func NewCommunicator(socket *websocket.Conn, syncFolder string) *Communicator {
	return &Communicator{
		conn:         socket,
		FileCreated:  make(chan sync.FileWithData),
		FilesDeleted: make(chan []string),
		syncFolder:   syncFolder,
	}
}

func isDisconnect(err error) bool {
	return strings.Contains(err.Error(), "closure")
}

func (c *Communicator) sendNewFileMessage(file *sync.FileWithData) error {
	event, err := fileCreatedEvent(file, c.syncFolder)
	if err != nil {
		return err
	}

	return c.conn.WriteJSON(event)
}

func (c *Communicator) handleNewFileMessage(e Event) {
	var file sync.FileWithData
	err := json.Unmarshal(e.Json, &file)
	if err != nil {
		log.Error(err)
		return
	}
	c.FileReceived(file)
}

// TODO: Handle events not strings
func (c *Communicator) listenToClient() {
	for {
		var e Event
		err := c.conn.ReadJSON(&e)
		if err != nil {
			if isDisconnect(err) {
				return
			}
			log.Errorf("Error reading json: %s", err)

			return
		}

		switch e.Id {
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
			// log.Debug("Communicator received file")
			c.sendNewFileMessage(&file)
		case deletedFiles := <-c.FilesDeleted:
			log.Debugf("Handling deleted files: %v", deletedFiles)
		}
	}
}

func (c *Communicator) HandleComms(conn *websocket.Conn) {
	c.conn = conn
	go c.writeToClient()

	c.listenToClient()
}
