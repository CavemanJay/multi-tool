package config

import (
	"encoding/json"
	"os"

	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("multi-tool")

type ServerOptions struct {
}

type ClientOptions struct {
	Host string
}

type ArchiveOptions struct {
	OutFolder string
	InFolders []string
}

type MusicOptions struct {
	SecretsFile  string
	PlaylistName string
	Limit        int
}

type Config struct {
	// ServerOptions  ServerOptions `json:"-"`
	ClientOptions
	ArchiveOptions
	MusicOptions
	UseLastRun    bool `json:"-"`
	Recursive     bool
	Port          int
	AppDataFolder string
	SyncFolder    string
	Append        string
	ConfigPath    string
}

func ReadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		// log.Errorf("Error reading config file: %s", err)
		return nil, err
	}
	defer file.Close()

	cfg := Config{}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&cfg)
	if err != nil {
		// log.Errorf("Error parsing config file: %s", err)
		return nil, err
	}

	return &cfg, nil
}

func WriteConfig(path string, cfg *Config) error {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		// log.Errorf("Error opening config file for writing: %s", err)
		return err
	}
	defer file.Close()
	// fmt.Printf("Writing to file: %s\n", path)

	file.Truncate(0)
	encoder := json.NewEncoder(file)
	err = encoder.Encode(cfg)
	if err != nil {
		log.Error(err)
	}

	return err
}
