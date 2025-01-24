/* Copyright (c) 2021 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package cli

import (
	"errors"
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/jievince/bubbline"
	"github.com/jievince/bubbline/editline"
	"github.com/vesoft-inc/nebula-ng-tools/ngql/pkg/completer"
)

type Editor interface {
	ReadInput(p string) (string, error)
	LoadHistory(historyFile string) error
	AddHistory(item string) error
	SaveHistory() error
	Close()
}

// BubblineEditor implements Editor interface using bubbline library
type BubblineEditor struct {
	editor *bubbline.Editor
}

// NewBubblineEditor creates a new editor with bubbline
func NewBubblineEditor(historyFile string) Editor {
	editor := bubbline.New()
	t := &BubblineEditor{editor: editor}

	// Configure editor
	editor.SetExternalEditorEnabled(true, "ngql")
	editor.NextPrompt = "-> "
	editor.AutoComplete = editline.AutoCompleteFn(t.getCompletions)
	editor.CheckInputComplete = t.checkInputComplete
	editor.SetDebugEnabled(true)

	// set static cursor style, avoid flickering
	editor.CursorMode = cursor.CursorStatic

	prevSiStyle := editor.FocusedStyle.SearchInput
	editor.FocusedStyle = editline.Style{SearchInput: prevSiStyle}
	editor.BlurredStyle = editline.Style{}

	err := editor.LoadHistory(historyFile)
	if err != nil {
		log.Panicf("Load history file %s failed, %s", historyFile, err.Error())
	}
	editor.SetAutoSaveHistory(historyFile, false)

	return t
}

// It processes multi-line input by removing triple quotes and backslashes
// while preserving the original line structure
func (t *BubblineEditor) ReadInput(p string) (string, error) {
	t.editor.Prompt = p
	line, err := t.editor.GetLine()
	if err != nil {
		if errors.Is(err, bubbline.ErrInterrupted) {
			return "", ErrPromptAborted
		}
		return "", err
	}

	lines := strings.Split(line, "\n")
	if len(lines) <= 1 {
		return line, nil
	}

	// Check if input is wrapped in triple quotes
	firstLine := strings.TrimLeft(lines[0], " \t")

	if strings.HasPrefix(firstLine, "'''") || strings.HasPrefix(firstLine, `"""`) {
		var quoteType string
		if strings.HasPrefix(firstLine, "'''") {
			quoteType = "'''"
		} else {
			quoteType = `"""`
		}

		// Process first line: remove triple quotes from start, keep the rest
		lines[0] = firstLine[3:]
		startLineIndex := 0
		if lines[0] == "" {
			startLineIndex = 1
		}

		// Find and process last line with matching triple quotes
		for i := len(lines) - 1; i >= 0; i-- {
			trimmedLine := strings.TrimRight(lines[i], " \t")
			if strings.HasSuffix(trimmedLine, quoteType) {
				// Remove triple quotes from end, keep the rest
				lines[i] = trimmedLine[:len(trimmedLine)-3]
				endLineIndex := i + 1
				if lines[i] == "" {
					endLineIndex = i
				}
				// Keep only lines between start and end quotes (inclusive)
				lines = lines[startLineIndex:endLineIndex]
				break
			}
		}

		return strings.Join(lines, "\n"), nil
	}

	// Handle backslash line continuation
	var result []string
	for i, l := range lines {
		if i < len(lines)-1 {
			// Find the last non-space character position
			lastNonSpace := len(l) - 1
			for lastNonSpace >= 0 && (l[lastNonSpace] == ' ' || l[lastNonSpace] == '\t') {
				lastNonSpace--
			}

			// Check if the last non-space character is a backslash
			if lastNonSpace >= 0 && l[lastNonSpace] == '\\' {
				// Keep everything before the backslash (preserving leading spaces)
				result = append(result, l[:lastNonSpace])
			} else {
				result = append(result, l)
			}
		} else {
			// Last line
			result = append(result, l)
		}
	}

	return strings.Join(result, "\n"), nil
}

// LoadHistory loads the entry history from file
func (t *BubblineEditor) LoadHistory(historyFile string) error {
	return t.editor.LoadHistory(historyFile)
}

// AddHistory adds a history entry and optionally saves
func (t *BubblineEditor) AddHistory(item string) error {
	return t.editor.AddHistory(item)
}

// SaveHistory saves the current history to the file
func (t *BubblineEditor) SaveHistory() error {
	return t.editor.SaveHistory()
}

// Close cleans up resources
func (t *BubblineEditor) Close() {
	t.editor.Close()
}

func (t *BubblineEditor) getCompletions(entireInput [][]rune, line, col int) (string, editline.Completions) {
	text := string(entireInput[line])
	_, completions, _ := completer.NewCompleter(text, col)
	return "", editline.SimpleWordsCompletion(completions, "completion", col, 0, len(text))
}

// checkInputComplete checks if the input is complete
func (t *BubblineEditor) checkInputComplete(v [][]rune, line, col int) bool {
	// Return true if input is empty
	if len(v) == 0 || (len(v) == 1 && len(v[0]) == 0) {
		return true
	}

	// Check triple quotes mode first
	if t.isInTripleQuotes(v) {
		return false
	}

	// Check if the last line ends with backslash (line continuation)
	lastLine := string(v[len(v)-1])
	return !strings.HasSuffix(strings.TrimRight(lastLine, " \t"), "\\")
}

// isInTripleQuotes checks if the input is within triple quotes
// Triple quotes (”' or """) are used to wrap multi-line input
// The opening triple quotes must be at the start of the first line (ignoring leading spaces)
// The closing triple quotes must be at the end of some line (ignoring trailing spaces)
func (t *BubblineEditor) isInTripleQuotes(v [][]rune) bool {
	if len(v) == 0 {
		return false
	}

	// Check if first line starts with triple quotes (ignoring leading spaces)
	firstLine := strings.TrimLeft(string(v[0]), " \t")
	if !strings.HasPrefix(firstLine, "'''") && !strings.HasPrefix(firstLine, `"""`) {
		return false
	}

	// Determine quote type
	quoteType := `"""`
	if strings.HasPrefix(firstLine, "'''") {
		quoteType = "'''"
	}

	// Check if any line ends with matching triple quotes
	for i := len(v) - 1; i > 0; i-- {
		line := strings.TrimRight(string(v[i]), " \t")
		if strings.HasSuffix(line, quoteType) {
			return false
		}
	}

	// No matching end quote found
	return true
}
