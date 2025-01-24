package cli

type Cli interface {
	ReadInput() (line string, exit bool, err error)
	GetPrompt() string
	Output() bool
	Close()
}
