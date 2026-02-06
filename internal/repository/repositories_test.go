package repository

import (
	"authentication/models"
	"os"
	"testing"
)

func TestCreateUser(t *testing.T) {
	temp, err := os.CreateTemp("/home/aida-najafi/Go/project2.0/testFiles", "users-*.csv")
	if err != nil {
		t.Fatal(err)
	}
	// defer os.Remove(temp.Name())
	repo := NewCSVRepo(temp.Name())
	user := models.User{
		ID:       1,
		Name:     "test",
		Username: "testUser",
		Email:    "test@example.com",
		Password: "hashed",
	}
	err = repo.Create(user)
	if err != nil {
		t.Fatalf("Test failed: %v", err)
	}
}


func TestGenerateID(t *testing.T){
	temp, err:= os.CreateTemp("/home/aida-najafi/Go/project2.0/testFiles", "users-*.csv")
	if err!= nil{
		t.Fatal(err)
	}
	// defer os.Remove(temp)
	repo:= NewCSVRepo(temp.Name())
	id1, _:= repo.GenerateID()
	id2, _ := repo.GenerateID()
	if id1 != 1{
		t.Fatalf("must have been 1, got %v", id1)
	}
	if id2 != 2{
		t.Fatalf("must have been 2, got %v", id2)
	}
}