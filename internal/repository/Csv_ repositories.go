package repository

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
)

type Repository struct {
	mu       sync.Mutex
	filePath string
}



func NewRepo(path string) *Repository {
	return &Repository{
		filePath: path,
	}
}

func (r *Repository) Init() error {
	file, err := os.OpenFile(r.filePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		if os.IsExist(err) {
			return nil
		}
		return fmt.Errorf("Failed to open file: %w", err)
	}

	defer file.Close()

	writer := csv.NewWriter(file)
	if err := writer.Write([]string{"id",
		"name",
		"username",
		"email",
		"password_hash"}); err != nil {
		return err
	}
	writer.Flush()
	return writer.Error()
}

func (r *Repository) CreateUser(user UserInfo) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	file, err := os.OpenFile(r.filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("Failed to open file: %w", err)
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()
	record := []string{
		strconv.Itoa(user.ID),
		user.Name,
		user.Username,
		user.Email,
		user.Password,
	}
	return writer.Write(record)
}

func (r *Repository) GenerateID() (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	file, err := os.Open(r.filePath)
	if err != nil {
		return 1, nil
	}
	defer file.Close()
	reader := csv.NewReader(file)
	lastID := 0
	lineCount := 0
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return 0, err
		}
		lineCount++
		if lineCount == 1 {
			continue
		}
		if id, err := strconv.Atoi(record[0]); err == nil {
			lastID = id
		}
	}
	if lastID == 0 {
		return 1, nil
	}
	return lastID + 1, nil
}

func (r *Repository) ReadAll() ([]UserInfo, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	file, err := os.Open(r.filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("Failed to read file path: %w", err)
	}
	var users []UserInfo
	for i, rec := range records {
		if i == 0 {
			continue
		}
		if len(rec) < 5 {
			continue
		}
		id, err := strconv.Atoi(rec[0])
		if err != nil {
			return nil, fmt.Errorf("Failed to convert id to int: %w", err)
		}
		user := UserInfo{
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

func (r *Repository) GetUserByUsername(username string) (*UserInfo, error) {
	users, err := r.ReadAll()
	if err != nil {
		return &UserInfo{}, fmt.Errorf("Failed to read path: %w", err)
	}
	for _, u := range users {
		if strings.TrimSpace(u.Username) == username {
			u.Name = strings.TrimSpace(u.Name)
			u.Username = strings.TrimSpace(u.Username)
			u.Email = strings.TrimSpace(u.Email)
			u.Password = strings.TrimSpace(u.Password)
			return &u, nil
		}

	}
	return &UserInfo{}, errors.New("User not found!")
}

func (r *Repository) AlreadyExistCheck(username, email string) error {

	users, err := r.ReadAll()
	if err != nil {
		return fmt.Errorf("Failed to read file: %w", err)
	}
	for _, u := range users {
		if strings.TrimSpace(username) == strings.TrimSpace(u.Username) || strings.TrimSpace(email) == strings.TrimSpace(u.Email) {

			return fmt.Errorf("user already exists : %s", username)
		}
	}
	return nil
}
