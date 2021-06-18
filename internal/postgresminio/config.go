package postgresminio

import (
	"crypto/rand"
	"math/big"
	"strings"
)

// Config contains all values needed for a working PostgreSQL storage service.
type Config struct {
	IDs struct {
		EntryIDLength   int
		RemovalIDLength int
		Chars           string
		charsRunes      []rune
	}
}

// Init must be called before interacting with the Config.
func (config *Config) Init() {
	config.IDs.charsRunes = []rune(config.IDs.Chars)
}

// GenerateEntryID generates a new entry ID based on the values read from the configuration struct.
func (config Config) GenerateEntryID() (string, error) {
	return config.generateUnfixedLengthID(config.IDs.EntryIDLength)
}

// GenerateRemovalID generates a new removal ID based on the values read from the configuration struct.
func (config Config) GenerateRemovalID() (string, error) {
	return config.generateUnfixedLengthID(config.IDs.RemovalIDLength)
}

// generateUnfixedLengthID is an internal function in order to generate IDs based on the configuration values and using
// the secure crypto/rand package.
func (config Config) generateUnfixedLengthID(length int) (string, error) {
	if config.IDs.charsRunes == nil {
		panic("config Init function has to be called first")
	}
	var stringBuilder strings.Builder
	int64CharsRunesLength := int64(len(config.IDs.charsRunes))
	for i := 0; i < length; i++ {
		bigInt, err := rand.Int(rand.Reader, big.NewInt(int64CharsRunesLength))
		if err != nil {
			return "", err
		}
		randomRune := config.IDs.charsRunes[int(bigInt.Int64())]
		_, err = stringBuilder.WriteRune(randomRune)
		if err != nil {
			return "", err
		}
	}
	return stringBuilder.String(), nil
}
