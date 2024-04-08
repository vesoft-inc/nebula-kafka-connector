/* Copyright (c) 2020 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package cli

import (
	"bufio"
	"fmt"
	"io"
)

// non-interactive
type nCli struct {
	status
	buf    *bufio.Reader
	io     io.ReadCloser
	output bool
}

func NewnCli(i io.ReadCloser, output bool, user string) Cli {
	return &nCli{
		status: status{
			user:                 user,
			space:                "(none)",
			respErr:              "",
			playingData:          false,
			promptLen:            -1,
			promptColor:          -1,
			line:                 "",
			joinedByTripleQuotes: false,
			joinedByBackSlash:    false,
		},
		io:     i,
		buf:    bufio.NewReader(i),
		output: output,
	}
}

func readln(r *bufio.Reader) (string, error) {
	var (
		isPartial bool  = true
		err       error = nil
		line, ln  []byte
	)
	for isPartial && err == nil {
		line, isPartial, err = r.ReadLine()
		ln = append(ln, line...)
	}
	return string(ln), err
}

func (l *nCli) Output() bool {
	return l.output
}

func (l *nCli) ReadLine() (string, bool, error) {
	for {
		input, err := readln(l.buf)
		if err == nil {
			if l.output {
				fmt.Print(l.status.nebulaPrompt())
				// not record input to historyFile now
				fmt.Println(input)
			}
			l.status.checkJoined(input)
			if l.status.joinedByTripleQuotes || l.status.joinedByBackSlash {
				continue
			}
			return l.status.line, false, nil
		} else if err == io.EOF {
			return "", true, nil
		} else {
			return "", false, err
		}
	}
}

func (l *nCli) Interactive() bool {
	return false
}

func (l *nCli) SetRespError(msg string) {
	l.status.respErr = msg
}

func (l *nCli) GetRespError() string {
	return l.status.respErr
}

func (l *nCli) SetSpace(space string) {
	if len(space) > 0 {
		l.status.space = space
	} else {
		l.status.space = "(none)"
	}
}

func (l *nCli) GetSpace() string {
	return l.status.space
}

func (l *nCli) PlayingData(b bool) {
	l.playingData = b
}

func (l nCli) IsPlayingData() bool {
	return l.playingData
}

func (l *nCli) Close() {
	l.io.Close()
}

func (l *nCli) GetPrompt() string {
	return l.status.nebulaPrompt()
}
