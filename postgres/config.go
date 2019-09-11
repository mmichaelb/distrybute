package postgres

import (
	"crypto/rand"
	"gopkg.in/yaml.v2"
	"io"
	"math/big"
	"os"
	"strings"
)

// Config contains all values needed for a working PostgreSQL storage service.
type Config struct {
	IDs struct {
		EntryIDLength   int    `yaml:"entry_id_length"`
		RemovalIDLength int    `yaml:"removal_id_length"`
		Chars           string `yaml:"chars"`
		charsRunes      []rune
	} `yaml:"entry_ids"`
}

// ReadConfigurationFile reads from the given filepath and if successful returns the parsed configuration.
func ReadConfigurationFile(filePath string) (*Config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := file.Close(); err !=nil {
			panic(err) // if this err is non-nil panic is the only way out
		}
	}()
	// forward the read process to the ReadConfiguration function
	return ReadConfiguration(file)
}

// ReadConfiguration reads from the given file and if successful returns the parsed configuration.
func ReadConfiguration(reader io.Reader) (*Config, error) {
	decoder := yaml.NewDecoder(reader)
	var config *Config
	// use yaml decoder to decode the configuration file
	err := decoder.Decode(config)
	return config, err
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
