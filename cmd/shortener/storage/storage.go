package storage

type Repository interface {
	GetData(string) (string, error)
	SaveData(string, string) error
}

var Store Repository
