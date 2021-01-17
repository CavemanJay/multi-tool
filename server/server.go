package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/JayCuevas/jays-server/sync"
	"github.com/gorilla/websocket"
)

type Server struct {
	fileCreated chan *sync.File
	upgrader    websocket.Upgrader
	newFiles    []*sync.File
	Port        int
	fileWatcher sync.FileWatcher
	// clients     map[net.Addr]*websocket.Conn
}

func NewServer(rootFolder string, recursive bool, port int) *Server {
	server := &Server{
		fileCreated: make(chan *sync.File),
		Port:        port,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin:     func(r *http.Request) bool { return true },
		},
	}

	server.fileWatcher = sync.FileWatcher{
		Root:        rootFolder,
		Recursive:   recursive,
		FileCreated: server.fileCreatedHandler,
	}

	return server
}

func (s *Server) Listen() error {
	exit := make(chan struct{}, 1)

	http.HandleFunc("/", s.socketEndpoint())
	go s.fileWatcher.Watch(exit)
	log.Printf("Server watching folder: %s", s.fileWatcher.Root)

	log.Printf("Server listening on :%d", s.Port)
	return http.ListenAndServe(fmt.Sprintf(":%d", s.Port), nil)
}

func (s *Server) sendFile(conn *websocket.Conn, file *sync.File) error {
	// log.Printf("Sending file to client: %s", file.Path)
	jsonData, err := file.ToJson()
	if err != nil {
		return err
	}
	return conn.WriteMessage(websocket.TextMessage, jsonData)
}

func (s *Server) writeToClient(conn *websocket.Conn) {
	// Send files that have been added before the client connected
	for _, file := range s.newFiles {
		s.sendFile(conn, file)
	}

	var file *sync.File
	for {
		file = <-s.fileCreated
		s.sendFile(conn, file)
	}
}

func (s *Server) listenToClient(conn *websocket.Conn) {
	for {
		// messageType, p, err := conn.ReadMessage()
		// if err != nil {
		// 	log.Println(err)
		// 	return
		// }
		// log.Println(string(p))
	}
}

func (s *Server) socketEndpoint() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ws, err := s.upgrader.Upgrade(w, r, nil)

		if err != nil {
			log.Println(err)
			return
		}

		defer ws.Close()

		log.Println("Client connected")
		go s.writeToClient(ws)
		s.listenToClient(ws)
	}
}

func (s *Server) fileCreatedHandler(file *sync.File) error {
	// Iterate over the files and see if the file has already been handled
	for i, f := range s.newFiles {
		if f.Path == file.Path {
			// If the hashes are the same then the file has been handled already
			if f.Hash == file.Hash {
				return nil
			}

			// The hashes are different, update the existing copy in memory
			s.newFiles[i] = file
		}
	}

	s.newFiles = append(s.newFiles, file)

	// Non blocking sending of the events
	select {
	case s.fileCreated <- file:
		// default:
	}

	return nil
}
