package commands

type Start struct {
}

func (c Start) Execute() string {
	return "I am Start Command"
}
