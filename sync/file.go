package sync

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"
)

type file struct {
	Path   string
	Hash   string
	Base64 string `json:"data"`
}

type File struct {
	Path string
	Hash string
}

func (f File) ToJson() ([]byte, error) {
	toSend, err := f.getFileToSend()
	if err != nil {
		return nil, err
	}
	return json.Marshal(toSend)
}

func (f File) getFileToSend() (*file, error) {
	encoded, err := getBase64FileData(f.Path)
	if err != nil {
		return nil, err
	}
	return &file{
		Path:   f.Path,
		Hash:   f.Hash,
		Base64: encoded,
	}, nil
}

func hashFile(path string) (string, error) {
	//Initialize variable returnMD5String now in case an error has to be returned
	var returnMD5String string

	file, err := os.Open(path)

	if err != nil {
		return returnMD5String, err
	}

	defer file.Close()

	//Open a new hash interface to write to
	hash := md5.New()

	//Copy the file in the hash interface and check for any error
	if _, err := io.Copy(hash, file); err != nil {
		return returnMD5String, err
	}

	//Get the 16 bytes hash
	hashInBytes := hash.Sum(nil)[:16]

	//Convert the bytes to a string
	returnMD5String = hex.EncodeToString(hashInBytes)

	return returnMD5String, nil
}

func readFile(path string) ([]byte, error) {
	file, err := os.Open(path)

	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func getBase64FileData(path string) (string, error) {

	data, err := readFile(path)

	if err != nil {
		log.Printf("Error reading file: %v", err)
		return "", err
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

func GetFileInfo(path string) (*File, error) {
	hash, err := hashFile(path)

	if err != nil {
		log.Printf("Error hashing file: %v", err)
		return nil, err
	}

	return &File{
		Hash: hash,
		Path: path,
	}, nil
}
