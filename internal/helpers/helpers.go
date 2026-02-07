package helpers

import "golang.org/x/crypto/bcrypt"

func HashPassword(pass string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(pass), 10)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func CheckPassword(plain, hashed string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(plain), []byte(hashed))
	return err == nil
}
