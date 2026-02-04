package repository

import (
	"encoding/csv"
	"errors"
	"os"
	"sync"
)

type CsvRepository struct {
	mu       sync.Mutex
	filePath string
}

func NewCSVRepo(path string) *CsvRepository {
	return &CsvRepository{
		filePath: path,
	}
}

func (r *CsvRepository) Init() error {
	_, err := os.Stat(r.filePath)
	if os.IsNotExist(err) {
		file, err := os.Create(r.filePath)
		if err != nil {
			return errors.New("File Doesn't exist and cannot be created!")
		}
		defer file.Close()
		writer := csv.NewWriter(file)
		defer writer.Flush()
		return writer.Write([]string{
			"id",
			"name",
			"username",
			"email",
			"password_hash",
		})
	}
	return nil

}
