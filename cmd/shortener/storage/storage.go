package storage

import (
	"context"
	"github.com/fngoc/url-shortener/internal/models"
)

type Repository interface {
	GetData(context.Context, string) (string, error)
	GetAllData(context.Context) ([]models.ResponseDto, error)
	DeleteData(userID int, url string) error
	SaveData(context.Context, string, string) error
}

var Store Repository
