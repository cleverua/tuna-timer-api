package commands

import (
	"fmt"
	"testing"
	"github.com/stretchr/testify/assert"
)


// let some code duplicate stay here...
func TestGetSimpleStartCommand(t *testing.T) {
	cmd, err := Get("start Convert the logotype to PNG")
	assert.NoError(t, err)

	commandType := fmt.Sprintf("%T", cmd)
	assert.Equal(t, "commands.Start", commandType)

	start := cmd.(Start)
	assert.Equal(t, "Convert the logotype to PNG", start.arguments["main"])
}

func TestGetSimpleStartCommandWithUnicodeArgument(t *testing.T) {
	cmd, err := Get("start Сконвертировать логотип в PNG")
	assert.NoError(t, err)

	commandType := fmt.Sprintf("%T", cmd)
	assert.Equal(t, "commands.Start", commandType)

	start := cmd.(Start)
	assert.Equal(t, "Сконвертировать логотип в PNG", start.arguments["main"])
}