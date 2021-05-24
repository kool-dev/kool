package environment

import (
	"os"

	"github.com/fireworkweb/godotenv"
)

// DefaultEnvStorage holds data to store environment variables
type DefaultEnvStorage struct{}

// EnvStorage contract that holds environment variables storage logic
type EnvStorage interface {
	Get(string) string
	Set(string, string)
	Load(string) error
	All() []string
	IsTrue(string) bool
}

// NewEnvStorage creates a new Environment Storage instance
func NewEnvStorage() EnvStorage {
	return &DefaultEnvStorage{}
}

// Get get environment variable value
func (es *DefaultEnvStorage) Get(key string) string {
	return os.Getenv(key)
}

// Set set environment variable value
func (es *DefaultEnvStorage) Set(key string, value string) {
	os.Setenv(key, value)
}

// Load load environment file
func (es *DefaultEnvStorage) Load(filename string) error {
	return godotenv.Load(filename)
}

// All get all environment variables
func (es *DefaultEnvStorage) All() []string {
	return os.Environ()
}

// IsTrue checks whether the given environment variable is
// to what would be a boolean value of true (either 1 or "true")
func (es *DefaultEnvStorage) IsTrue(key string) bool {
	value := os.Getenv(key)
	return value == "1" || value == "true"
}
