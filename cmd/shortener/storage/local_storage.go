package storage

import (
	"context"
	"fmt"
	"github.com/fngoc/url-shortener/internal/logger"
)

type LocalStore map[string]string

var localStorage LocalStore

func InitializeInMemoryLocalStore() error {
	localStorage = make(map[string]string)
	Store = localStorage
	logger.Log.Info("Initializing local storage")
	return nil
}

func (lc LocalStore) GetData(_ context.Context, key string) (string, error) {
	if val, ok := lc[key]; ok {
		return val, nil
	}
	return "", fmt.Errorf("data by key: %s, not found", key)
}

func (lc LocalStore) SaveData(_ context.Context, key string, value string) error {
	if key == "" || value == "" {
		return fmt.Errorf("key or value is empty")
	}
	if _, ok := lc[key]; ok {
		return fmt.Errorf("data by key: %s, already exists", key)
	}
	lc[key] = value
	return nil
}
