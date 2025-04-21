package kvstore

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// DefaultMaxFileSize is 64MB
const DefaultMaxFileSize = 64

// NewStore creates and initializes a new Store
func NewStore(dataDir string) (*Store, error) {
	// Create data directory if it doesn't exist
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	store := &Store{
		DataDir:     dataDir,
		MaxFileSize: DefaultMaxFileSize,
		KeyDir:      make(map[string]KeyLocation),
		Mutex:       sync.RWMutex{},
	}

	// Create initial data file
	if err := store.rotateActiveFile(); err != nil {
		return nil, fmt.Errorf("failed to create initial data file: %w", err)
	}

	return store, nil
}

// rotateActiveFile creates a new data file and makes it the active file
func (s *Store) rotateActiveFile() error {
	// Close the current active file if it exists
	if s.ActiveFile != nil {
		if err := s.ActiveFile.Close(); err != nil {
			return fmt.Errorf("failed to close active file: %w", err)
		}
	}

	// Generate a new filename with timestamp
	timestamp := time.Now().UnixNano()
	filename := fmt.Sprintf("data_%d.db", timestamp)
	filepath := filepath.Join(s.DataDir, filename)

	// Create and open the new file
	file, err := os.OpenFile(filepath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("failed to create new data file: %w", err)
	}

	s.ActiveFile = file
	s.ActiveFilename = filename
	s.ActiveFileSize = 0

	return nil
}

// Put writes a key-value pair to the store
func (s *Store) Put(key, value string) error {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	entry := Entry{
		TimeStamp: time.Now().Unix(),
		Key:       key,
		Value:     value,
		KeySize:   int32(len(key)),
		ValueSize: int32(len(value)),
	}

	// Check if we need to rotate the file
	entrySize := int64(16 + len(key) + len(value)) // timestamp(8) + keySize(4) + valueSize(4) + key + value
	if s.ActiveFileSize+entrySize > s.MaxFileSize {
		if err := s.rotateActiveFile(); err != nil {
			return fmt.Errorf("failed to rotate file: %w", err)
		}
	}

	// Write the entry to the active file
	fullPath := filepath.Join(s.DataDir, s.ActiveFilename)
	offset, err := WriteEntry(fullPath, entry)
	if err != nil {
		return fmt.Errorf("failed to write entry: %w", err)
	}

	// Update the KeyDir with the new location
	s.KeyDir[key] = KeyLocation{
		Filename: s.ActiveFilename,
		Offset:   offset,
	}

	// Update the active file size
	s.ActiveFileSize += entrySize

	return nil
}

// Get retrieves a value by key from the store
func (s *Store) Get(key string) (string, error) {
	s.Mutex.RLock()
	defer s.Mutex.RUnlock()

	// Look up the key location
	location, exists := s.KeyDir[key]
	if !exists {
		return "", fmt.Errorf("key not found: %s", key)
	}

	// Read the entry from the appropriate file
	fullPath := filepath.Join(s.DataDir, location.Filename)
	value, err := ReadEntry(fullPath, location.Offset)
	if err != nil {
		return "", fmt.Errorf("failed to read entry: %w", err)
	}

	return value, nil
}

// Close closes the active file
func (s *Store) Close() error {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	if s.ActiveFile != nil {
		return s.ActiveFile.Close()
	}
	return nil
}
