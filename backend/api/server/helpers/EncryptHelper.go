package helpers

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

/*Hashing Password*/
func Hashing(pass string) string {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(pass), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(hashedBytes)
}

/*Match password for login*/
func VerifyPassword(enteredPassword string, storedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(enteredPassword))

	return err == nil
}
