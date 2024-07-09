package storage

import "fmt"

type Repository interface {
	GetData(string) (string, error)
	SaveData(string, string) error
}

type LocalStore map[string]string

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
	return nil
}

var Store = make(LocalStore)
