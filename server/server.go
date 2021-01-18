package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/JayCuevas/jays-server/sync"
	"github.com/gorilla/websocket"
)

type Server struct {
	fileCreated chan *sync.File
	upgrader    websocket.Upgrader
	filesList   []*sync.File
	Port        int
	fileWatcher sync.FileWatcher
	clientCount int
}

func NewServer(rootFolder string, recursive bool, port int) *Server {
	server := &Server{
		fileCreated: make(chan *sync.File),
		Port:        port,
		filesList:   []*sync.File{},
		clientCount: 0,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin:     func(r *http.Request) bool { return true },
		},
	}

	server.fileWatcher = sync.FileWatcher{
		Root:         rootFolder,
		Recursive:    recursive,
		FileCreated:  server.fileCreatedHandler,
		Files:        &server.filesList,
		FilesDeleted: server.filesDeletedHandler,
	}

	return server
}

func (s *Server) Listen() error {
	exit := make(chan struct{}, 1)
	// termChan := make(chan os.Signal)
	// signal.Notify(termChan, syscall.SIGTERM, syscall.SIGINT)

	http.HandleFunc("/ws", s.socketEndpoint())
	http.HandleFunc("/files", s.filesEndpoint())

	log.Printf("Server watching folder: %s", s.fileWatcher.Root)
	log.Printf("Server listening on :%d", s.Port)

	go func() {
		err := s.fileWatcher.IndexFiles(func(file *sync.File) {
			s.filesList = append(s.filesList, file)
		})

		if err != nil {
			log.Printf("Error indexing files: %v", err)
			return
		}
	}()
	// go func() {
	// 	<-termChan

	// }()
	go s.fileWatcher.Watch(exit)
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
	// for _, file := range s.filesList {
	// 	s.sendFile(conn, file)
	// }

	var file *sync.File
	for {
		file = <-s.fileCreated
		s.sendFile(conn, file)
	}
}

func (s *Server) listenToClient(conn *websocket.Conn) {
	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		log.Println(string(p))
	}
}

func (s *Server) filesEndpoint() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(s.filesList)
	}
}

func (s *Server) socketEndpoint() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ws, err := s.upgrader.Upgrade(w, r, nil)

		if err != nil {
			log.Println(err)
			return
		}

		defer func() {
			ws.Close()
			s.clientCount--
			log.Printf("Client disconnected: %s", ws.RemoteAddr().String())
		}()

		s.clientCount++
		log.Printf("Client connected: %s", ws.RemoteAddr().String())

		go s.writeToClient(ws)
		s.listenToClient(ws)
	}
}

func (s *Server) fileCreatedHandler(file *sync.File) error {
	// Iterate over the files and see if the file has already been handled
	for i, f := range s.filesList {
		if f.Path == file.Path {
			// If the hashes are the same then the file has been handled already
			if f.Hash == file.Hash {
				return nil
			}

			// The hashes are different, update the existing copy in memory
			s.filesList[i] = file
		}
	}

	s.filesList = append(s.filesList, file)

	if s.clientCount > 0 {
		// This will block if there are no clients connected
		select {
		case s.fileCreated <- file:
			// default:
		}
	}

	return nil
}

func removeFile(files []*sync.File, i int) []*sync.File {
	files[i] = files[len(files)-1]
	return files[:len(files)-1]
}

func (s *Server) filesDeletedHandler(paths []string) error {
	for _, path := range paths {
		for i, file := range s.filesList {
			if path == file.Path {
				s.filesList = removeFile(s.filesList, i)
			}
		}
	}

	return nil
}
