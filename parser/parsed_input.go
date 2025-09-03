package parser

type ParsedInput struct {
	Command     string
	Subcommands []string
	InputFlags  []InputFlag
}

func NewParserInput(command string) *ParsedInput {
	return &ParsedInput{
		Command:     command,
		Subcommands: make([]string, 0),
		InputFlags:  make([]InputFlag, 0),
	}
}
