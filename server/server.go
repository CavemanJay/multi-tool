package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	comms "github.com/CavemanJay/gogurt/communications"
	"github.com/CavemanJay/gogurt/database"
	filesync "github.com/CavemanJay/gogurt/sync"
	"github.com/gorilla/websocket"
	"github.com/op/go-logging"
)

var (
	lock sync.Mutex
	log  = logging.MustGetLogger("gogurt")
)

type Server struct {
	upgrader     websocket.Upgrader
	filesList    []filesync.File
	Port         int
	fileWatcher  filesync.FileWatcher
	clientCount  int
	communicator *comms.Communicator
	db           *database.Manager
}

func NewServer(rootFolder string, recursive bool, port int) *Server {
	db, err := database.NewManager()
	if err != nil {
		log.Fatalf("Error initializing server: %s", err)
	}
	err = db.ApplyMigrations()
	if err != nil {
		log.Fatalf("Error applying database migrations: %s", err)
	}

	server := &Server{
		Port:         port,
		filesList:    []filesync.File{},
		clientCount:  0,
		communicator: comms.NewCommunicator(nil, rootFolder),
		db:           db,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin:     func(r *http.Request) bool { return true },
		},
	}

	server.fileWatcher = filesync.FileWatcher{
		Root:        rootFolder,
		Recursive:   recursive,
		FileCreated: server.fileCreatedHandler,
		Files:       &server.filesList,
		FileDeleted: server.fileDeletedHandler,
	}

	return server
}

func (s *Server) addFile(file filesync.File) {
	// lock.Lock()
	// s.filesList = append(s.filesList, file)
	s.db.Upsert(&file)
	// lock.Unlock()
}

func (s *Server) Listen() error {
	defer s.db.Close()
	exit := make(chan struct{}, 1)

	http.HandleFunc("/ws", s.socketEndpoint())
	http.HandleFunc("/files", s.filesEndpoint())

	log.Debugf("Server watching folder: %s", s.fileWatcher.Root)
	log.Infof("Server listening on localhost:%d", s.Port)

	go func() {
		err := s.fileWatcher.IndexFiles(s.addFile)

		if err != nil {
			log.Errorf("Error indexing files: %v", err)
			return
		}
		log.Debug("Done indexing files")
	}()
	go s.fileWatcher.Watch(exit)
	return http.ListenAndServe(fmt.Sprintf(":%d", s.Port), nil)
}

func (s *Server) filesEndpoint() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(s.db.All())
	}
}

func (s *Server) handleClientDisconnect(ws *websocket.Conn) {
	ws.Close()
	s.clientCount--
	log.Infof("Client disconnected: %s", ws.RemoteAddr().String())
	log.Debugf("%d client(s) remaining", s.clientCount)
}

func (s *Server) socketEndpoint() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ws, err := s.upgrader.Upgrade(w, r, nil)

		if err != nil {
			log.Error(err)
			return
		}

		defer s.handleClientDisconnect(ws)

		s.clientCount++
		log.Infof("Client connected: %s", ws.RemoteAddr().String())

		s.communicator.HandleComms(ws)
	}
}

func (s *Server) fileCreatedHandler(file filesync.File) error {
	s.addFile(file)

	if s.clientCount > 0 {
		f, err := file.ToDataFile(s.fileWatcher.Root)
		if err != nil {
			return err
		}
		// This will block if there are no clients connected
		s.communicator.FileCreated <- f
	}

	return nil
}

// TODO
func (s *Server) fileDeletedHandler(path string) error {
	// s.db.Delete(path)

	return nil
}
