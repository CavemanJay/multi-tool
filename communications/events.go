package communications

import (
	"encoding/json"

	"github.com/CavemanJay/gogurt/sync"
)

type eventId int

const (
	EventFileCreated eventId = iota
	EventFileDeleted
)

type Event struct {
	Id   eventId
	Json []byte
}

func fileCreatedEvent(file *sync.FileWithData, root string) (*Event, error) {
	json, err := json.Marshal(file)
	if err != nil {
		return nil, err
	}

	return &Event{
		Id:   EventFileCreated,
		Json: json,
	}, nil
}
