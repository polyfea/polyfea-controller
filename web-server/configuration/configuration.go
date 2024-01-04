package configuration

import (
	"os"
	"strings"
)

var configuration = map[string]string{}

// GetConfigurationValueOrDefault returns the value of the configuration key.
// If the configuration key is not found, it returns default value.
// First call will load the configuration from different sources.
// Supported sources are: Environment variables.
// The sources are loaded in the same order as listed above.
func GetConfigurationValueOrDefault(key string, defaultValue string) string {
	if len(configuration) == 0 {
		initConfiguration()
	}

	if len(configuration[key]) > 0 {
		return configuration[key]
	}

	return defaultValue
}

// GetConfigurationValue returns the value of the configuration key.
// If the configuration key is not found, it returns an empty string.
// First call will load the configuration from different sources.
// Supported sources are: Environment variables.
// The sources are loaded in the same order as listed above.
func GetConfigurationValue(key string) string {
	return GetConfigurationValueOrDefault(key, "")
}

func initConfiguration() {
	// Load configuration from environment variables
	for _, env := range os.Environ() {
		configuration[strings.Split(env, "=")[0]] = strings.Split(env, "=")[1]
	}
}
