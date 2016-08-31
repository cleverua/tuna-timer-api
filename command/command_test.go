package command

import (
	"testing"
)

func TestGet(t *testing.T) {
	_, err := Get("unknown")

	if err == nil {
		t.Error("Expected to get an error because command not found")
	}
}
