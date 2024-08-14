package storage

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/fngoc/url-shortener/internal/logger"
	"github.com/fngoc/url-shortener/internal/models"
	"os"
)

type FileStore map[string]string

var (
	fileStorage     FileStore
	currentUUID     int
	currentFilePath string
)

func InitializeFileLocalStore(filename string) error {
	logger.Log.Info("Initializing file store")
	currentFilePath = filename

	if ok, _ := isCreate(currentFilePath); !ok {
		_, err := os.Create(currentFilePath)
		if err != nil {
			return err
		}
		fileStorage = make(FileStore)
		Store = fileStorage
		return nil
	}

	file, err := os.OpenFile(currentFilePath, os.O_RDONLY, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	fileStorage = make(FileStore)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		saveData := models.URLData{}
		data := scanner.Bytes()
		err := json.Unmarshal(data, &saveData)
		if err != nil {
			return err
		}
		currentUUID = saveData.UUID
		fileStorage[saveData.ShortURL] = saveData.OriginalURL
	}
	Store = fileStorage
	return nil
}

func (fs FileStore) GetData(_ context.Context, key string) (string, error) {
	if val, ok := fs[key]; ok {
		return val, nil
	}
	return "", fmt.Errorf("data by key: %s, not found", key)
}

func (fs FileStore) SaveData(_ context.Context, key string, value string) error {
	if key == "" || value == "" {
		return fmt.Errorf("key or value is empty")
	}
	if _, ok := fs[key]; ok {
		return fmt.Errorf("data by key: %s, already exists", key)
	}
	fs[key] = value

	err := saveToFile(key, value)
	if err != nil {
		return err
	}

	return nil
}

func isCreate(filename string) (bool, error) {
	file, err := os.Open(filename)
	if err != nil {
		return false, err
	}
	defer file.Close()
	return true, nil
}

func saveToFile(shortURL, OriginalURL string) error {
	currentUUID += 1
	saveData := models.URLData{UUID: currentUUID, ShortURL: shortURL, OriginalURL: OriginalURL}
	data, err := json.Marshal(saveData)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(currentFilePath, os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	_, err = file.WriteString(string(data) + "\n")
	return err
}
