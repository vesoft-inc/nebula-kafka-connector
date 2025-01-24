package cli

import (
	"errors"
	"fmt"
	"io"
)

// ErrPromptAborted is the error for prompt aborted. i.e. ctrl+d.
var ErrPromptAborted = errors.New("prompt aborted")

// interactive
type iCli struct {
	editor Editor
	user   string
}

func NewiCli(historyFile, user string) Cli {
	t := NewBubblineEditor(historyFile)

	return &iCli{
		editor: t,
		user:   user,
	}
}

func (l *iCli) ReadInput() (string, bool, error) {
	input, err := l.editor.ReadInput(l.GetPrompt())
	if err != nil {
		if err == ErrPromptAborted {
			return "", false, nil
		} else if err == io.EOF {
			return "", true, nil
		}
		return "", false, err
	}

	if len(input) > 0 {
		l.editor.AddHistory(input)
	}
	return input, false, nil
}

func (l *iCli) GetPrompt() string {
	return fmt.Sprintf("(%s@nebula)> ", l.user)
}

func (l *iCli) Output() bool {
	return true
}

func (l *iCli) Close() {
	defer l.editor.Close()
	l.editor.SaveHistory()
}
