package config

import (
	"errors"
	"strconv"
	"strings"
)

// Repository represents a contract to handle configuration reading
type Repository interface {
	GetSystemConfiguration() (Configuration, error)
}

// Configuration represents a list of configurantion for many environments
type Configuration struct {
	CurrentEnvironment string
	Environments       []EnvironmentConfig
}

// EnvironmentConfig represents an environment item configuration
type EnvironmentConfig struct {
	Name       string
	HTTPServer HTTPServerConfig
}

// HTTPServerConfig represents the HTTP Server address configuration
type HTTPServerConfig struct {
	Address        string
	Port           int
	RequestTimeout int
}

// GetFormatedAddress returns a formated HTTP Address
func (h HTTPServerConfig) GetFormatedAddress() string {
	return strings.Join([]string{h.Address, strconv.Itoa(h.Port)}, ":")
}

// GetCurrentEnvironmentConfig returns the current environment configuration
func (c Configuration) GetCurrentEnvironmentConfig() (EnvironmentConfig, error) {
	for _, environmentConfig := range c.Environments {
		if environmentConfig.Name == c.CurrentEnvironment {
			return environmentConfig, nil
		}
	}

	return EnvironmentConfig{}, errors.New("Environment name not found")
}
