package commands

// Start - starts timer for specific task
// If there is an other started task then it will be stopped
type Start struct {
}

func (c Start) Execute() CommandResult {
	return CommandResult{}
}
