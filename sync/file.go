package sync

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type FileWithData struct {
	Path   string
	Hash   string
	Base64 string `json:"data"`
}

type File struct {
	Path string
	Hash string
}

func (f File) ToDataFile(syncRoot string) (FileWithData, error) {
	toSend, err := f.getFileToSend(syncRoot)
	return *toSend, err
}

func (f File) getFileToSend(root string) (*FileWithData, error) {
	encoded, err := getBase64FileData(filepath.Join(root, f.Path))
	if err != nil {
		return nil, err
	}
	return &FileWithData{
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
		log.Errorf("Error reading file: %v", err)
		return "", err
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

func getFilePathRelativeToRoot(root string, fullPath string) string {
	root, _ = filepath.Abs(root)
	fullPath, _ = filepath.Abs(fullPath)

	path := strings.ReplaceAll(fullPath, root, "")
	path = strings.ReplaceAll(path, "\\", "/")
	if path[0] == '/' {
		path = path[1:]
	}

	return path
}

func GetFileInfo(syncRoot string, path string) (*File, error) {
	hash, err := hashFile(path)

	if err != nil {
		log.Errorf("Error hashing file: %v", err)
		return nil, err
	}

	return &File{
		Hash: hash,
		Path: getFilePathRelativeToRoot(syncRoot, path),
	}, nil
}
