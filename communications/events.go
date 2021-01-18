package communications

import "github.com/JayCuevas/jays-server/sync"

type eventId int

const (
	EventFileCreated eventId = iota
	EventFileDeleted
)

type event struct {
	id   eventId
	json []byte
}

func fileCreatedEvent(file *sync.File) (*event, error) {
	data, err := file.ToJson()
	if err != nil {
		return nil, err
	}
	return &event{
		id:   EventFileCreated,
		json: data,
	}, nil
}
