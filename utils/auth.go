package utils

import (
	"crypto/sha256"
	"fmt"

	"github.com/google/uuid"
)

const SESSION_COOKIE_KEY = "sid"

// FIXME: Очень плохой способ хранения паролей, но это не является частью текущей задачи, нужно хотябы использовать соль
func HashPassword(pwd string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(pwd)))
}

func CheckPassword(pwd string, hashed_pwd string) bool {
	return HashPassword(pwd) == hashed_pwd
}

func CreateSID() string {
	return uuid.New().String()

}
