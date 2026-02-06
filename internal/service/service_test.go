package service

import (
	"authentication/internal/repository"
	"os"
	"testing"
)

func TestSignUp(t *testing.T) {
	temp, err := os.CreateTemp("/home/aida-najafi/Go/project2.0/testFiles", "users-*.csv")
	if err != nil {
		t.Fatal(err)
	}
	// defer os.Remove(temp.Name())
	repo := repository.NewCSVRepo(temp.Name())
	service := NewAuthService(repo)
	info := SignUpInfo{
		Name:     "test",
		Username: "testUser",
		Email:    "test@example.com",
		Password: "hashed",
	}
	user, err := service.SignUp(info)
	if err != nil {
		t.Fatal(err)
	}
	if user.Password == "hashed" {
		t.Fatal("password was not hashed")
	}

	if user.ID == 0 {
		t.Fatal("generateID failed")
	}

}
