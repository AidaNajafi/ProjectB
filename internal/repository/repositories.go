package repository

import (
	"authentication/internal/helpers"
	"authentication/models"
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strconv"
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

func (r *CsvRepository) Create(user models.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	file, err := os.OpenFile(r.filePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return errors.New("Failed to open file")
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()
	hasehdPass, _ := helpers.HashPassword(user.Password)
	record := []string{
		strconv.Itoa(user.ID),
		user.Name,
		user.Username,
		user.Email,
		hasehdPass,
	}
	return writer.Write(record)
}

func (r *CsvRepository) GenerateID() (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	file, err := os.Open(r.filePath)
	if err != nil {
		return 1, nil
	}
	defer file.Close()
	reader := csv.NewReader(file)
	record, err := reader.ReadAll()
	if err != nil {
		return 0, err
	}
	if len(record) <= 1 {
		return 1, nil
	}
	lastRow := record[len(record)-1]
	id, err := strconv.Atoi(lastRow[0])
	if err != nil {
		return 0, err
	}
	return id + 1, nil
}

func (r *CsvRepository) ReadAll() ([]models.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	file, err := os.Open(r.filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file", err)
	}
	defer file.Close()
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("Failed to read file path", err)
	}
	var users []models.User
	for _, rec := range records {
		id, err := strconv.Atoi(rec[0])
		if err != nil {
			return nil, fmt.Errorf("Failed to convert id to int", err)
		}
		user := models.User{
			ID:       id,
			Name:     rec[1],
			Username: rec[2],
			Email:    rec[3],
			Password: rec[4],
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *CsvRepository) GetUserByUsername(username string) (models.User, error) {
	users, err := r.ReadAll()
	if err != nil {
		return models.User{}, fmt.Errorf("Failed to read path", err)
	}
	for _, u := range users {
		if u.Username == username {
			return u, nil
		}
	}
	return models.User{}, errors.New("User not found!")
}
