package commands

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestGetUnknownCommand(t *testing.T) {
	_, err := Get("unknown")
	assert.Error(t, err)
}
