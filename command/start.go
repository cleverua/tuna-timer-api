package command

type StartCommand struct {
}

func (c StartCommand) Execute() string {
	return "I am Start Command"
}
