package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetEnv(t *testing.T) {
	key := "TEST_GET_ENV_UNIQUE_KEY"

	testCases := []struct {
		name         string
		setValue     string
		defaultValue string
		want         string
	}{
		{
			name:         "unset var returns default",
			setValue:     "",
			defaultValue: "fallback",
			want:         "fallback",
		}, {
			name:         "set var returns its value",
			setValue:     "hello",
			defaultValue: "fallback",
			want:         "hello",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setValue != "" {
				t.Setenv(key, tc.setValue)
			}
			got := getEnv(key, tc.defaultValue)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestGetEnvAsInt(t *testing.T) {
	key := "TEST_GET_ENV_INT_UNIQUE_KEY"

	testCases := []struct {
		name         string
		setValue     string
		defaultValue int
		want         int
	}{
		{
			name:         "valid integer string is parsed",
			setValue:     "8080",
			defaultValue: 0000,
			want:         8080,
		},
		{
			name:         "invalid integer string fallback to default",
			setValue:     "notInteger",
			defaultValue: 0000,
			want:         0000,
		},
		{
			name:         " unset var returns default",
			setValue:     "",
			defaultValue: 0000,
			want:         0000,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setValue != "" {
				t.Setenv(key, tc.setValue)
			}
			got := getEnvAsInt(key, tc.defaultValue)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestGetEnvAsDuration(t *testing.T) {
	key := "TEST_GET_ENV_DURATION_UNIQUE_KEY"

	testCases := []struct {
		name         string
		setValue     string
		defaultValue time.Duration
		want         time.Duration
	}{
		{
			name:         "valid duration string is parsed",
			setValue:     "30s",
			defaultValue: 15 * time.Second,
			want:         30 * time.Second,
		},
		{
			name:         "invalid duration string fallback to default",
			setValue:     "notDuration",
			defaultValue: 15 * time.Second,
			want:         15 * time.Second,
		},
		{
			name:         "unset var returns default",
			setValue:     "",
			defaultValue: 15 * time.Second,
			want:         15 * time.Second,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setValue != "" {
				t.Setenv(key, tc.setValue)
			}
			got := getEnvAsDuration(key, tc.defaultValue)
			assert.Equal(t, tc.want, got)
		})
	}
}
