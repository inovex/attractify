package auth

import (
	"bytes"
	"crypto/rand"

	"golang.org/x/crypto/argon2"
)

type Password struct {
	Password []byte
	Salt     []byte
}

func NewPassword(password string) Password {
	pw := Password{}
	pw.Salt = randGen(16)
	pw.Password = hashPassword(password, pw.Salt)

	return pw
}

func (p Password) Compare(password string) bool {
	pwd := hashPassword(password, p.Salt)
	return bytes.Equal(pwd, p.Password)
}

func randGen(size int) []byte {
	buf := make([]byte, size)
	if _, err := rand.Read(buf); err != nil {
		return nil
	}
	return buf
}

func hashPassword(password string, salt []byte) []byte {
	return argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 64)
}
