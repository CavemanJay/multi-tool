package main

import (
	"log"
	"net/http"

	"github.com/JayCuevas/jays-server/sync"
	"github.com/gorilla/websocket"
)

var (
	fileCreated = make(chan *sync.File)

	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}

	newFiles []*sync.File
)

func reader(conn *websocket.Conn) {
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		log.Println(string(p))

		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}
	}
}

func sendFile(conn *websocket.Conn, file *sync.File) error {
	log.Printf("Sending file to client: %s", file.Path)
	return conn.WriteJSON(file)
}

func writer(conn *websocket.Conn) {
	// Send files that have been added before the client connected
	for _, file := range newFiles {
		sendFile(conn, file)
	}

	var file *sync.File
	for {
		file = <-fileCreated
		sendFile(conn, file)
	}
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err)
		return
	}

	defer ws.Close()

	log.Println("Client connected")
	go writer(ws)
	reader(ws)
}

func setupRoutes() {
	http.HandleFunc("/", wsEndpoint)
}

func fileCreatedHandler(file *sync.File) error {
	// Iterate over the files and see if the file has already been handled
	for i, f := range newFiles {
		if f.Path == file.Path {
			// If the hashes are the same then the file has been handled already
			if f.Hash == file.Hash {
				return nil
			}

			// The hashes are different, update the existing copy in memory
			newFiles[i] = file
		}
	}

	newFiles = append(newFiles, file)

	// for _, f := range newFiles {
	// 	log.Printf("%v", f.Path)
	// }

	// Non blocking sending of the events
	select {
	case fileCreated <- file:
	default:
	}

	return nil
}

func main() {
	setupRoutes()

	exit := make(chan struct{}, 1)

	root := "/mnt/c/Users/cueva/Sync"
	recursive := true

	fileWatcher := sync.FileWatcher{Root: root, Recursive: recursive, FileCreated: fileCreatedHandler}

	go fileWatcher.Watch(exit)
	// fileWatcher.Watch(exit)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
