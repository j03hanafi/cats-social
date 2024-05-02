package security

import (
	"golang.org/x/crypto/bcrypt"

	"cats-social/common/configs"
)

// HashPassword generates a secure hash of a password using bcrypt.
// Returns the hashed password as a hex-encoded string or an error if hashing fails.
// example for making salt - https://play.golang.org/p/_Aw6WeWC42I
func HashPassword(password string) (string, error) {
	// using recommended cost parameters from - https://godoc.org/golang.org/x/crypto/scrypt
	bhash, err := bcrypt.GenerateFromPassword([]byte(password), configs.Runtime.API.BCryptSalt)
	if err != nil {
		return "", err
	}

	// return hex-encoded string with salt appended to password
	hashedPW := string(bhash)

	return hashedPW, nil
}

// ComparePasswords checks if the supplied password matches the stored hashed password.
// Returns true if they match or an error if there's a problem in the comparison process.
func ComparePasswords(storedPassword, suppliedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(suppliedPassword))
}
