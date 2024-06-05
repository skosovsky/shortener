package store

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"shortener/internal/log"
	"shortener/internal/service"
)

type FileStore struct {
	file    *os.File
	encoder *json.Encoder
	*MemoryStore
}

func NewFileStore(fileName string) (*FileStore, error) {
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("open file %s error: %w", fileName, err)
	}

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("get file %s stat error: %w", fileName, err)
	}

	memoryStore := NewMemoryStore()

	if fileInfo.Size() != 0 {
		scanner := bufio.NewScanner(file)

		for scanner.Scan() {
			data := scanner.Bytes()

			var site service.Site
			if err = json.Unmarshal(data, &site); err != nil {
				return nil, fmt.Errorf("unmarshal json error: %w", err)
			}

			_ = memoryStore.Add(site) // err nil
		}

		if err = scanner.Err(); err != nil {
			return nil, fmt.Errorf("scanner error: %w", err)
		}
	}

	return &FileStore{
		file:        file,
		encoder:     json.NewEncoder(file),
		MemoryStore: memoryStore,
	}, nil
}

func (f *FileStore) Add(site service.Site) error {
	_ = f.MemoryStore.Add(site) // err nil

	err := f.encoder.Encode(site)
	if err != nil {
		return fmt.Errorf("encode file %s error: %w", f.file.Name(), err)
	}

	return nil
}

func (f *FileStore) Close() {
	err := f.file.Close()
	if err != nil {
		log.Error("close file %s error", f.file.Name(),
			log.ErrAttr(err))
	}
}
