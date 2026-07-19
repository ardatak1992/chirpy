package auth

import (
	"runtime"

	"github.com/alexedwards/argon2id"
)

// HashPassword hashes a password using argon2id and returns the hash string.
func HashPassword(password string) (string, error) {

	params := argon2id.Params{
		SaltLength:  16,
		KeyLength:   32,
		Iterations:  3,
		Parallelism: uint8(runtime.NumCPU()),
		Memory:      68000,
	}

	hash, err := argon2id.CreateHash(password, &params)
	if err != nil {
		return "", err
	}

	return hash, nil

}

// CheckPasswordHash compares password and hashed password and returns a boolean
func CheckPasswordHash(password, hash string) (bool, error) {
	return argon2id.ComparePasswordAndHash(password, hash)
}
