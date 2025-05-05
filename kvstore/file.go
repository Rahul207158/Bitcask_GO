package kvstore

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

func WriteEntry(filepath string, entry Entry) (int64, error) {
	file, err := os.OpenFile(filepath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return 0, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Get current offset
	offset, err := file.Seek(0, io.SeekEnd)
	if err != nil {
		return 0, fmt.Errorf("failed to seek to end: %w", err)
	}

	// Write timestamp
	if err = binary.Write(file, binary.LittleEndian, entry.TimeStamp); err != nil {
		return 0, fmt.Errorf("failed to write timestamp: %w", err)
	}

	// Write key size
	if err = binary.Write(file, binary.LittleEndian, entry.KeySize); err != nil {
		return 0, fmt.Errorf("failed to write key size: %w", err)
	}

	// Write value size
	if err = binary.Write(file, binary.LittleEndian, entry.ValueSize); err != nil {
		return 0, fmt.Errorf("failed to write value size: %w", err)
	}

	// Write entry type
	if err = binary.Write(file, binary.LittleEndian, entry.Type); err != nil {
		return 0, fmt.Errorf("failed to write entry type: %w", err)
	}

	// Write key
	if _, err = file.Write([]byte(entry.Key)); err != nil {
		return 0, fmt.Errorf("failed to write key: %w", err)
	}

	// Write value
	if _, err = file.Write([]byte(entry.Value)); err != nil {
		return 0, fmt.Errorf("failed to write value: %w", err)
	}

	return offset, nil
}

func ReadEntry(filepath string, offset int64) (string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Skip timestamp (8 bytes)
	_, err = file.Seek(offset+8, io.SeekStart)
	if err != nil {
		return "", fmt.Errorf("failed to seek to entry: %w", err)
	}

	// Read keySize, valueSize, and type (4 + 4 + 1 = 9 bytes)
	sizeBuf := make([]byte, 9)
	if _, err := io.ReadFull(file, sizeBuf); err != nil {
		return "", fmt.Errorf("failed to read size buffer: %w", err)
	}
	keySize := int32(binary.LittleEndian.Uint32(sizeBuf[0:4]))
	valueSize := int32(binary.LittleEndian.Uint32(sizeBuf[4:8]))
	entryType := EntryType(sizeBuf[8])

	// If this is a tombstone entry, return an error
	if entryType == EntryTombstone {
		return "", fmt.Errorf("key has been deleted")
	}

	// Skip the key
	_, err = file.Seek(int64(keySize), io.SeekCurrent)
	if err != nil {
		return "", fmt.Errorf("failed to skip key: %w", err)
	}

	// Read the value
	valueBuf := make([]byte, valueSize)
	if _, err := io.ReadFull(file, valueBuf); err != nil {
		return "", fmt.Errorf("failed to read value: %w", err)
	}

	return string(valueBuf), nil
}
