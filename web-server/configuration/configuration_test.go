package configuration

import (
	"os"
	"testing"
)

func TestGetConfigurationValueValueReturned(t *testing.T) {
	// Arrange
	testKey := "TEST_KEY"
	testValue := "TEST_VALUE"
	os.Setenv(testKey, testValue)

	// Act
	value := GetConfigurationValue(testKey)

	// Assert
	if value != testValue {
		t.Errorf("GetConfigurationValue(%s) = %s; want %s", testKey, value, testValue)
	}
}

func TestGetConfigurationValueEmptyStringReturnedWhenKeyNotFound(t *testing.T) {
	// Arrange
	testKey := "TEST_KEY_EMPTY"

	// Act
	value := GetConfigurationValue(testKey)

	// Assert
	if value != "" {
		t.Errorf("GetConfigurationValue(%s) = %s; want %s", testKey, value, "")
	}
}

func TestGetConfigurationValueDefaultValueReturnedWhenKeyNotFound(t *testing.T) {
	// Arrange
	testKey := "TEST_KEY_DEFAULT_VALUE"
	testDefaultValue := "TEST_DEFAULT_VALUE"

	// Act
	value := GetConfigurationValueOrDefault(testKey, testDefaultValue)

	// Assert
	if value != testDefaultValue {
		t.Errorf("GetConfigurationValueOrDefault(%s, %s) = %s; want %s", testKey, testDefaultValue, value, testDefaultValue)
	}
}
