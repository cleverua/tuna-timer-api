package config

import (
	//_ "github.com/jinzhu/gorm/dialects/postgres"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestFoo(t *testing.T) {
	SetupEnvironment("test")
	assert.Equal(t, 1, 1)
}
