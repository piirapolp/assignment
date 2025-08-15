package util

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"golang.org/x/crypto/sha3"
	"math/big"
	"strings"
)

func HashPassword(plainTextPassword string) (string, error) {
	salt := GenerateSalt()
	return HashPasswordFixSalt(salt, plainTextPassword)
}

func GenerateSalt() string {
	charSet := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	randomString, _ := GenerateRandomStringFromSpecificCharacters(charSet, 64)
	return randomString
}

func GenerateRandomStringFromSpecificCharacters(characters string, length int) (string, error) {
	if len(characters) == 0 || length <= 0 {
		return "", errors.New("invalid argument size")
	}

	randomString := ""
	for i := 0; i < length; i++ {
		value, _ := rand.Int(rand.Reader, big.NewInt(int64(len(characters))))
		randomString += string(characters[value.Int64()])
	}
	return randomString, nil
}

func HashPasswordFixSalt(salt, plainTextPassword string) (string, error) {
	buffer := []byte(plainTextPassword)

	hmac := make([]byte, 32)
	hash := sha3.NewShake256()

	if _, err := hash.Write([]byte(salt)); err != nil {
		return "", err
	}
	if _, err := hash.Write(buffer); err != nil {
		return "", err
	}
	if _, err := hash.Read(hmac); err != nil {
		return "", err
	}
	return salt + ":" + hex.EncodeToString(hmac), nil
}

func ValidatePin(plainPin, hashedPin string) (bool, error) {
	//get salt from hashedPassword
	splitHashedPin := strings.Split(hashedPin, ":")
	if len(splitHashedPin) != 2 {
		return false, fmt.Errorf("wrong hashed pin format to check")
	}
	salt := splitHashedPin[0]
	hashedFromPlain, err := HashPasswordFixSalt(salt, plainPin)
	if err != nil {
		return false, err
	}

	return hashedFromPlain == hashedPin, nil
}

func GenerateTokenSessionId(userId string) string {
	randomBytes, _ := GenerateRandomBytes(32)
	hash := sha256.New()
	hash.Write([]byte(fmt.Sprintf("%s", userId)))
	hash.Write(randomBytes)
	sum := hash.Sum(nil)
	return hex.EncodeToString(sum)
}

func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}
