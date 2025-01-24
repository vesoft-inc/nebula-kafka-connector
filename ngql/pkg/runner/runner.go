package runner

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/mattn/go-runewidth"
	nebula "github.com/vesoft-inc/nebula-ng-tools/golang"
	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/errors"
	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/types"
	"github.com/vesoft-inc/nebula-ng-tools/ngql/pkg/cli"
	"github.com/vesoft-inc/nebula-ng-tools/ngql/pkg/printer"
)

const (
	default_pager        = true
	default_pagerLimit   = 200
	default_pagerCommand = "less"
	maxRetryTimes        = 1
)

type (
	Runner struct {
		option       *runnerOption
		stdout       io.Writer
		file         io.WriteCloser
		client       types.Client
		sessionId    int64
		cli          cli.Cli
		printer      printer.Printer
		combinOutput io.Writer
		running      bool
	}

	combinOutput struct {
		runner *Runner
	}

	runnerOption struct {
		interactive    bool
		enableOutput   bool
		fileReader     io.ReadCloser // for non-interactive mode
		address        string
		user           string
		password       string
		historyDir     string
		timeoutSec     int
		failFast       bool // if true, stop loop for error
		widthMax       int
		signalChan     <-chan os.Signal
		pager          bool
		pagerLimit     int
		pagerCommand   string
		enableTLS      bool
		ca             string
		cert           string
		key            string
		peerNameVerify bool
		peerName       string
	}

	runnerOptionsFn func(*runnerOption)
)

func NewRunner(opts ...runnerOptionsFn) (*Runner, error) {
	var runOpts runnerOption
	// test less command
	cmd := exec.Command(default_pagerCommand)
	if err := cmd.Run(); err != nil {
		runOpts = runnerOption{
			pager: false,
		}
	} else {
		runOpts = runnerOption{
			pager:        default_pager,
			pagerLimit:   default_pagerLimit,
			pagerCommand: default_pagerCommand,
		}
	}
	r := &Runner{
		stdout: os.Stdout,
		option: &runOpts,
	}
	r.combinOutput = &combinOutput{
		runner: r,
	}
	for _, opt := range opts {
		opt(r.option)
	}

	r.printer = printer.NewPrinter("default", r.option.widthMax)
	if !r.option.enableOutput {
		r.stdout = io.Discard
	}
	if err := r.openNebulaClient(); err != nil {
		return nil, err
	}
	var c cli.Cli
	if r.option.interactive {
		historyFile := filepath.Join(r.option.historyDir, ".nebula_history")
		// Create history file if not exists with proper permissions
		if _, err := os.Stat(historyFile); os.IsNotExist(err) {
			file, err := os.Create(historyFile)
			if err != nil {
				return nil, fmt.Errorf("failed to create history file: %w", err)
			}
			file.Close()
		}
		c = cli.NewiCli(historyFile, r.option.user)
	} else {
		c = cli.NewnCli(r.option.fileReader, r.option.enableOutput, r.option.user)
	}
	r.cli = c

	return r, nil
}

func WithInteractive(i bool) runnerOptionsFn {
	return func(o *runnerOption) {
		o.interactive = i
	}
}

func WithTimeoutSec(sec int) runnerOptionsFn {
	return func(o *runnerOption) {
		o.timeoutSec = sec
	}
}

func WithHistoryDir(dir string) runnerOptionsFn {
	return func(o *runnerOption) {
		o.historyDir = dir
	}
}

func WithReadCloser(r io.ReadCloser) runnerOptionsFn {
	return func(o *runnerOption) {
		o.fileReader = r
	}
}

func WithWidthMax(width int) runnerOptionsFn {
	return func(o *runnerOption) {
		o.widthMax = width
	}
}

func WithNebula(addr, user, password string) runnerOptionsFn {
	return func(o *runnerOption) {
		o.address = addr
		o.user = user
		o.password = password
	}
}

func WithTLS(enable bool, ca, cert, key string, peerNameVerify bool, peerName string) runnerOptionsFn {
	return func(o *runnerOption) {
		o.enableTLS = enable
		o.ca = ca
		o.cert = cert
		o.key = key
		o.peerNameVerify = peerNameVerify
		o.peerName = peerName
	}
}

func WithOutput(output bool) runnerOptionsFn {
	return func(o *runnerOption) {
		o.enableOutput = output
	}
}

func WithFailFast(failFast bool) runnerOptionsFn {
	return func(o *runnerOption) {
		o.failFast = failFast
	}
}

func WithSignalChan(signalChan <-chan os.Signal) runnerOptionsFn {
	return func(o *runnerOption) {
		o.signalChan = signalChan
	}
}

func (r *Runner) welcome() {
	if !r.option.interactive {
		return
	}
	r.printStdout("Welcome to NebulaGraph 5.0, the distributed graph database offering native GQL support!\n")
	r.printStdout(":help for help.\n")
	r.printStdout("\n")
}

func (r *Runner) bye() {
	r.printStdout(fmt.Sprintf("Bye %s!\n\n", r.option.user))
	r.printStdout(fmt.Sprintf("%s\n", time.Now().In(time.Local).Format(time.RFC1123)))
}

// Loop the request util fatal or timeout
// Add input break yourself as `SHOW \<CR>HOSTS`
func (r *Runner) loop() error {
	for {
		input, exit, err := r.cli.ReadInput()
		if err != nil {
			return err
		}
		if exit { // Ctrl+D
			r.printBoth("\n")
			return nil
		}
		if len(input) == 0 { // 1). The input is empty, or 2). user presses ctrlC so the input is truncated
			continue
		}
		input = strings.TrimSpace(input)
		// record in file
		r.printFile(fmt.Sprintf("%s%s\n", r.cli.GetPrompt(), input))
		r.running = true
		exit, err = r.execute(input)
		r.running = false
		if err != nil {
			if r.option.failFast {
				return err
			}
			// if session is not found, or connection error, should exit
			// otherwise do not exit, and continue to execute next statement
			ne, ok := err.(*errors.NebulaError)
			if ok {
				switch ne.Code() {
				case errors.ERROR_SESSION_NOT_FOUND:
					fallthrough
				case errors.ERROR_CONN_IS_BROKEN, errors.ERROR_CONN_IS_CLOSED:
					return err
				}
			}
			r.printBoth(fmt.Sprintf("Error: %s\n", err.Error()))
			// "42" means SYNTAX_ERROR_OR_ACCESS_RULE_VIOLATION error class
			if ok && ne.ErrorClass() == "42" {
				r.highlightSyntaxError(input, err.Error())
			}
			continue
		}
		if exit { // :exit
			return nil
		}
	}
}

// execute one line
func (r *Runner) execute(input string) (exit bool, err error) {
	cmd, err := getCommand(r, input)
	if err != nil {
		return false, err
	}
	if cmd != nil {
		if err := cmd.execute(); err != nil {
			return false, err
		}
		if cmd.cmd == commandExit || cmd.cmd == commandQuit {
			return true, nil
		} else {
			return false, nil
		}
	}
	// execute nebula statement
	isVertical := strings.HasSuffix(input, "\\G")
	if isVertical {
		input = strings.TrimSuffix(input, "\\G")
	}
	start := time.Now()
	resp, err := r.executeGQLWithRetry(input, 0, maxRetryTimes)
	if err != nil {
		if r.option.failFast {
			return true, err
		}

		// if session is not found, or connection error, should exit
		// otherwise do not exit, and continue to execute next statement
		ne, ok := err.(*errors.NebulaError)
		if ok {
			switch ne.Code() {
			case errors.ERROR_SESSION_NOT_FOUND:
				fallthrough
			case errors.ERROR_CONN_IS_BROKEN, errors.ERROR_CONN_IS_CLOSED:
				return true, err
			}
		}
	}
	duration := time.Since(start)
	if resp == nil {
		return false, err
	}
	var out io.WriteCloser
	var execCmd *exec.Cmd
	if r.option.pager && resp.RowSize() >= r.option.pagerLimit {
		execCmd = exec.Command(r.option.pagerCommand)
		out, err = execCmd.StdinPipe()
		if err != nil {
			return false, err
		}
		execCmd.Stdout = os.Stdout
		if err := execCmd.Start(); err != nil {
			return false, err
		}
		r.stdout = out
	}
	if isVertical {
		r.printer.PrintResultVertical(r.combinOutput, resp)
	} else {
		r.printer.PrintResult(r.combinOutput, resp)
	}
	if r.option.pager && resp.RowSize() >= r.option.pagerLimit {
		out.Close()
		execCmd.Wait()
		r.stdout = os.Stdout
	}
	lantencyInMs := resp.Summary().TotalServerTimeUs() * 1000
	if err != nil {
		r.printBoth(fmt.Sprintf("Execution failed (time spent %v/%v)\n\n",
			time.Duration(lantencyInMs), duration),
		)
	} else {
		if resp.RowSize() > 0 {
			numRows := resp.RowSize()
			r.printBoth(fmt.Sprintf("Got %d rows (time spent %v/%v)\n\n",
				numRows, time.Duration(lantencyInMs), duration),
			)
		} else {
			r.printBoth(fmt.Sprintf("Execution succeeded (time spent %v/%v)\n\n",
				time.Duration(lantencyInMs), duration),
			)
		}
		if resp.Summary().ExplainType() != "" {
			r.printer.PrintPlanInfo(r.combinOutput, resp.Summary())
		}
	}
	r.printBoth(fmt.Sprintf("%s\n\n", time.Now().In(time.Local).Format(time.RFC1123)))

	return false, err
}

func (r *Runner) executeGQLWithRetry(stmt string, retryTimes, retryLimit int) (types.Result, error) {
	if r.client == nil {
		if err := r.openNebulaClient(); err != nil {
			return nil, err
		}
	}
	resp, err := r.client.Execute(stmt)
	if err == nil {
		return resp, nil
	}
	ne, ok := err.(*errors.NebulaError)
	if !ok {
		return nil, err
	}
	if retryTimes+1 > retryLimit {
		return nil, err
	}
	switch ne.Code() {
	// retry for session not found
	case errors.ERROR_SESSION_NOT_FOUND:
		fallthrough

	// retry for connection error
	case errors.ERROR_CONN_IS_BROKEN, errors.ERROR_CONN_IS_CLOSED:
		if err := r.reopenNebulaClient(); err != nil {
			return nil, err
		}
		return r.executeGQLWithRetry(stmt, retryTimes+1, retryLimit)
	default:
		return nil, err
	}
}

func (r *Runner) reopenNebulaClient() error {
	if r.client != nil {
		r.client.Close()
	}
	return r.openNebulaClient()
}

func (r *Runner) openNebulaClient() error {
	var err error
	r.client, err = nebula.NewNebulaClient(r.option.address, r.option.user, r.option.password,
		nebula.WithClientRequestTimeout(time.Duration(r.option.timeoutSec)*time.Second),
		nebula.WithClientTLS(r.option.enableTLS, r.option.ca, r.option.cert, r.option.key, r.option.peerNameVerify, r.option.peerName),
	)
	if err != nil {
		return err
	}
	id, err := r.client.GetSessionId()
	if err != nil {
		return err
	}
	r.sessionId = id
	if err := r.client.Ping(); err != nil {
		return err
	}
	return nil
}

func (r *Runner) Run() error {
	if r.option.interactive {
		r.welcome()
		defer r.bye()
	}
	loopErr := make(chan error)
	killed := make(chan struct{})
	// when receive ctrl+C, if the runner is executing, would kill query
	go func() {
		for {
			select {
			case signal := <-r.option.signalChan:
				switch signal {
				case os.Interrupt:
					if r.running {
						if err := r.killQuery(); err != nil {
							r.printBoth(fmt.Sprintf("Error: %s\n", err.Error()))
						}
					}
				case syscall.SIGTERM:
					if r.running {
						if err := r.killQuery(); err != nil {
							r.printBoth(fmt.Sprintf("Error: %s\n", err.Error()))
						}
					}
					killed <- struct{}{}
				}
			}
		}
	}()
	go func() {
		if err := r.loop(); err != nil {
			loopErr <- err
		} else {
			loopErr <- nil
		}
	}()
	select {
	case <-killed:
		return nil
	case err := <-loopErr:
		return err
	}
}

func (r *Runner) Close() {
	if r.client != nil {
		r.client.Close()
	}
	if r.cli != nil {
		r.cli.Close()
	}
	if r.file != nil {
		r.file.Close()
	}
}

func (r *Runner) printStdout(s string) {
	fmt.Fprint(r.stdout, s)
}

func (r *Runner) printFile(s string) {
	if r.file != nil {
		fmt.Fprint(r.file, s)
	}
}

func (r *Runner) printBoth(s string) {
	r.printFile(s)
	r.printStdout(s)
}

func (r *Runner) killQuery() error {
	killSession, err := nebula.NewNebulaClient(r.option.address, r.option.user, r.option.password,
		nebula.WithClientRequestTimeout(time.Duration(r.option.timeoutSec)*time.Second),
		nebula.WithClientTLS(r.option.enableTLS, r.option.ca, r.option.cert, r.option.key, r.option.peerNameVerify, r.option.peerName),
	)
	if err != nil {
		return err
	}
	defer killSession.Close()
	showQueryStmt := fmt.Sprintf("call show_queries() filter where session_id = %d return query_id", r.sessionId)
	resp, err := killSession.Execute(showQueryStmt)
	if err != nil {
		return err
	}
	// sometimes the query is not in query list, should wait for a while
	if resp.RowSize() == 0 {
		return fmt.Errorf("no query to kill")
	}
	row, err := resp.Next()
	if err != nil {
		return err
	}
	query_id, err := row.GetValueByIndex(0)
	if err != nil {
		return err
	}
	killQueryStmt := fmt.Sprintf("call kill_query(\"%s\") return 1", query_id.String())
	resp, err = killSession.Execute(killQueryStmt)
	if err != nil {
		return err
	}
	return nil
}

func (o *combinOutput) Write(p []byte) (n int, err error) {
	n, err = o.runner.stdout.Write(p)
	if o.runner.file != nil {
		o.runner.file.Write(p)
	}
	return
}

// calculate the display width for error marker
func getRuneDisplayWidth(r rune) int {
	// Handle tab character
	if r == '\t' {
		return 8
	}
	// Handle control characters
	if r == '\r' || r == '\n' || r == '\b' || r == '\f' || r == '\v' {
		return 0
	}
	return runewidth.RuneWidth(r)
}

func (r *Runner) highlightSyntaxError(input string, errMsg string) {
	// Check if error message contains position information
	idx := strings.LastIndex(errMsg, "at [L")
	if idx == -1 {
		return
	}

	// Extract position information [L4:1-L5:3]
	loc := regexp.MustCompile(`L(\d+):(\d+)-L(\d+):(\d+)\]$`).FindStringSubmatch(errMsg)
	if len(loc) != 5 {
		return
	}

	// Parse position information
	startLineNum, _ := strconv.Atoi(loc[1])
	startColNum, _ := strconv.Atoi(loc[2])
	endLineNum, _ := strconv.Atoi(loc[3])
	endColNum, _ := strconv.Atoi(loc[4])

	// Split input into lines, handling different line endings (\n, \r\n, \r)
	input = strings.ReplaceAll(input, "\r\n", "\n")
	input = strings.ReplaceAll(input, "\r", "\n")
	lines := strings.Split(input, "\n")

	// Validate line numbers
	if startLineNum <= 0 || startLineNum > len(lines) || endLineNum <= 0 || endLineNum > len(lines) {
		return
	}

	r.printBoth("\n") // Add a blank line before error display

	// Print error lines with markers
	for lineNum := startLineNum; lineNum <= endLineNum; lineNum++ {
		line := lines[lineNum-1]
		runes := []rune(line)
		marker := &strings.Builder{}

		// Print line number and separator
		r.printBoth(fmt.Sprintf("%4d | %s\n", lineNum, line))
		r.printBoth("     | ")

		if lineNum == startLineNum && lineNum == endLineNum {
			// Single line case
			// Add spaces before the error marker
			for i := 0; i < startColNum-1; i++ {
				marker.WriteRune(' ')
			}
			// Calculate total width of characters in error range
			totalWidth := 0
			for i := startColNum - 1; i < endColNum && i < len(runes); i++ {
				totalWidth += getRuneDisplayWidth(runes[i])
			}
			// Add markers according to total width
			for i := 0; i < totalWidth; i++ {
				marker.WriteRune('^')
			}
		} else if lineNum == startLineNum {
			// Start line case
			for i := 0; i < startColNum-1; i++ {
				marker.WriteRune(' ')
			}
			// Calculate total width from start to end of line
			totalWidth := 0
			for i := startColNum - 1; i < len(runes); i++ {
				totalWidth += getRuneDisplayWidth(runes[i])
			}
			for i := 0; i < totalWidth; i++ {
				marker.WriteRune('^')
			}
		} else if lineNum == endLineNum {
			// End line case
			totalWidth := 0
			for i := 0; i < endColNum && i < len(runes); i++ {
				totalWidth += getRuneDisplayWidth(runes[i])
			}
			for i := 0; i < totalWidth; i++ {
				marker.WriteRune('^')
			}
		} else {
			// Middle lines case
			totalWidth := 0
			for _, r := range runes {
				totalWidth += getRuneDisplayWidth(r)
			}
			for i := 0; i < totalWidth; i++ {
				marker.WriteRune('^')
			}
		}

		r.printBoth(fmt.Sprintf("%s\n", marker.String()))
	}

	r.printBoth("\n")
}
