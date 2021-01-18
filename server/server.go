package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	comms "github.com/JayCuevas/gogurt/communications"
	filesync "github.com/JayCuevas/gogurt/sync"
	"github.com/gorilla/websocket"
)

var lock sync.Mutex

type Server struct {
	fileCreated  chan filesync.File
	upgrader     websocket.Upgrader
	filesList    []filesync.File
	Port         int
	fileWatcher  filesync.FileWatcher
	clientCount  int
	communicator *comms.Communicator
}

func NewServer(rootFolder string, recursive bool, port int) *Server {
	server := &Server{
		fileCreated: make(chan filesync.File),
		Port:        port,
		filesList:   []filesync.File{},
		clientCount: 0,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin:     func(r *http.Request) bool { return true },
		},
		communicator: comms.NewCommunicator(nil),
	}

	server.fileWatcher = filesync.FileWatcher{
		Root:         rootFolder,
		Recursive:    recursive,
		FileCreated:  server.fileCreatedHandler,
		Files:        &server.filesList,
		FilesDeleted: server.filesDeletedHandler,
	}

	return server
}

func (s *Server) addFile(file filesync.File) {
	lock.Lock()
	s.filesList = append(s.filesList, file)
	lock.Unlock()
	log.Println(len(s.filesList))
}

func (s *Server) Listen() error {
	exit := make(chan struct{}, 1)

	http.HandleFunc("/ws", s.socketEndpoint())
	http.HandleFunc("/files", s.filesEndpoint())

	log.Printf("Server watching folder: %s", s.fileWatcher.Root)
	log.Printf("Server listening on :%d", s.Port)

	go func() {
		err := s.fileWatcher.IndexFiles(s.addFile)

		if err != nil {
			log.Printf("Error indexing files: %v", err)
			return
		}
	}()
	go s.fileWatcher.Watch(exit)
	return http.ListenAndServe(fmt.Sprintf(":%d", s.Port), nil)
}

func (s *Server) filesEndpoint() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(s.filesList)
		log.Println(len(s.filesList))
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

		s.communicator.HandleComms(ws)
	}
}

func (s *Server) fileCreatedHandler(file filesync.File) error {
	// Iterate over the files and see if the file has already been handled
	for i, f := range s.filesList {
		if f.Path == file.Path {
			// If the hashes are the same then the file has been handled already
			if f.Hash == file.Hash {
				return nil
			}

			// The hashes are different, update the existing copy in memory
			lock.Lock()
			s.filesList[i] = file
			lock.Unlock()
		}
	}

	s.addFile(file)

	if s.clientCount > 0 {
		// This will block if there are no clients connected
		select {
		case s.fileCreated <- file:
			// default:
		}
	}

	return nil
}

func removeFile(files []*filesync.File, i int) []*filesync.File {
	lock.Lock()
	defer lock.Unlock()

	files[i] = files[len(files)-1]
	return files[:len(files)-1]
}

func (s *Server) filesDeletedHandler(paths []string) error {
	// for _, path := range paths {
	// 	for i, file := range s.filesList {
	// 		if path == file.Path {
	// 			s.filesList = removeFile(s.filesList, i)
	// 		}
	// 	}
	// }

	return nil
}
