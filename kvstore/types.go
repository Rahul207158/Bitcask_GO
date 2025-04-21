package kvstore

import (
	"os"
	"sync"
)

type Entry struct {
	TimeStamp int64
	KeySize   int32
	ValueSize int32
	Key       string
	Value     string
}

// FileInfo tracks information about a data file
type FileInfo struct {
	Filename string
	Size     int64
}

// Store manages the database files and operations
type Store struct {
	DataDir        string
	ActiveFile     *os.File
	ActiveFilename string
	ActiveFileSize int64
	MaxFileSize    int64
	KeyDir         map[string]KeyLocation
	Mutex          sync.RWMutex
}

// KeyLocation stores where a key's data can be found
type KeyLocation struct {
	Filename string
	Offset   int64
}

var KeyDir map[string]int64

type RequestPayload struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
