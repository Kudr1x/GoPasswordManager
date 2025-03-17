package crypto

import (
	"crypto/rand"
	"crypto/sha512"
	"golang.org/x/crypto/pbkdf2"
)

const (
	Iterations = 100000 // Фиксированное количество итераций
)

func GenerateSalt(size int) []byte {
	salt := make([]byte, size)
	rand.Read(salt)
	return salt
}

func DeriveKey(password string, salt []byte) []byte {
	return pbkdf2.Key(
		[]byte(password),
		salt,
		Iterations,
		64,
		sha512.New,
	)
}
