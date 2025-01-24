package cli

import (
	"fmt"
	"io"
	"strings"
)

// non-interactive
type nCli struct {
	reader BatchReader
	output bool
	user   string
}

func NewnCli(i io.ReadCloser, output bool, user string) Cli {
	return &nCli{
		reader: NewBatchReader(i),
		output: output,
		user:   user,
	}
}

func (l *nCli) ReadInput() (string, bool, error) {
	result, isEOF, _, err := l.reader.ReadInput()
	if err != nil {
		return "", false, err
	}

	if l.output && result != "" {
		// Print prompt and input only when there is actual content
		l.printPromptAndInput(result)
	}

	return result, isEOF, nil
}

func (l *nCli) GetPrompt() string {
	return fmt.Sprintf("(%s@nebula)> ", l.user)
}

func (l *nCli) Output() bool {
	return l.output
}

func (l *nCli) Close() {
	l.reader.Close()
}

func (l *nCli) printPromptAndInput(input string) {
	prompt := l.GetPrompt()
	// If the input contains multiple lines, add the prompt to each line
	lines := strings.Split(input, "\n")
	for i, line := range lines {
		if i == 0 {
			fmt.Print(prompt)
		} else {
			fmt.Print("-> ")
		}
		fmt.Println(line)
	}
}
