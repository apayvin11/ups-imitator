package model

import (
	"testing"
	"time"

)

func TestConfig(t *testing.T) *Config {
	return &Config{
		UpsSyncInterval: time.Second,
	}
}
