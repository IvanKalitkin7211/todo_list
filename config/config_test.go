package config

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	os.Setenv("DB_USER", "tester")
	defer os.Unsetenv("DB_USER")

	cfg := NewConfig()

	assert.NotNil(t, cfg)
}
