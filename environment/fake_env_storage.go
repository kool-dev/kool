package environment

import (
	"fmt"

	"github.com/fireworkweb/godotenv"
)

// FakeEnvStorage holds fake environment variables
type FakeEnvStorage struct {
	Envs        map[string]string
	CalledLoad  bool
	EnvsHistory map[string][]string
}

// NewFakeEnvStorage creates a new FakeEnvStorage
func NewFakeEnvStorage() *FakeEnvStorage {
	return &FakeEnvStorage{
		Envs:        make(map[string]string),
		EnvsHistory: make(map[string][]string),
	}
}

// Get get environment variable value (fake behavior)
func (f *FakeEnvStorage) Get(key string) string {
	return f.Envs[key]
}

// Set set environment variable value (fake behavior)
func (f *FakeEnvStorage) Set(key string, value string) {
	f.Envs[key] = value
	f.EnvsHistory[key] = append(f.EnvsHistory[key], value)
}

// Load load environment file (fake behavior)
func (f *FakeEnvStorage) Load(filename string) error {
	f.CalledLoad = true
	envs, _ := godotenv.Read(filename)
	for k, v := range envs {
		if _, exists := f.Envs[k]; !exists {
			f.Envs[k] = v
		}
	}
	return nil
}

// All get all environment variables
func (f *FakeEnvStorage) All() (envs []string) {
	for key, value := range f.Envs {
		envs = append(envs, fmt.Sprintf("%s=%s", key, value))
	}

	return
}

// IsTrue checks whether the given environment variable is
// to what would be a boolean value of true (fake behavior)
func (f *FakeEnvStorage) IsTrue(key string) bool {
	value := f.Envs[key]
	return value == "1" || value == "true"
}
