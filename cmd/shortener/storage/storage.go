package storage

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/fngoc/url-shortener/internal/models"
	"os"
)

type Repository interface {
	GetData(string) (string, error)
	SaveData(string, string) error
}

type LocalStore map[string]string

var Store LocalStore

var currentUUID int

var currentFilePath string

func InitializeLocalStore(filename string) error {
	currentFilePath = filename

	if ok, _ := isCreate(currentFilePath); !ok {
		_, err := os.Create(currentFilePath)
		if err != nil {
			return err
		}
		Store = make(LocalStore)
		return nil
	}

	file, err := os.OpenFile(currentFilePath, os.O_RDONLY, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	Store = make(LocalStore)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		saveData := models.URLData{}
		data := scanner.Bytes()
		err := json.Unmarshal(data, &saveData)
		if err != nil {
			return err
		}
		currentUUID = saveData.UUID
		Store[saveData.ShortURL] = saveData.OriginalURL
	}
	return nil
}

func (lc LocalStore) GetData(key string) (string, error) {
	if val, ok := lc[key]; ok {
		return val, nil
	}
	return "", fmt.Errorf("data by key: %s, not found", key)
}

func (lc LocalStore) SaveData(key string, value string) error {
	if key == "" || value == "" {
		return fmt.Errorf("key or value is empty")
	}
	if _, ok := lc[key]; ok {
		return fmt.Errorf("data by key: %s, already exists", key)
	}
	lc[key] = value

	err := saveToFile(key, value)
	if err != nil {
		return err
	}
	return nil
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

func isCreate(filename string) (bool, error) {
	file, err := os.Open(filename)
	if err != nil {
		return false, err
	}
	defer file.Close()
	return true, nil
}
