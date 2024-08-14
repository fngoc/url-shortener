package storage

import "context"

type Repository interface {
	GetData(context.Context, string) (string, error)
	SaveData(context.Context, string, string) error
}

var Store Repository
