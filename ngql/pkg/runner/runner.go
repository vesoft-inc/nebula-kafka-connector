package runner

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	nebula "github.com/vesoft-inc/nebula-ng-tools/golang"
	"github.com/vesoft-inc/nebula-ng-tools/ngql/pkg/cli"
	"github.com/vesoft-inc/nebula-ng-tools/ngql/pkg/printer"
)

var default_pager = true
var default_pagerLimit = 200
var default_pagerCommand = "less"

type (
	Runner struct {
		option       *runnerOption
		stdout       io.Writer
		file         io.WriteCloser
		client       nebula.Client
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
		enableGoPrompt bool
		timeoutSec     int
		failFast       bool //if true, stop loop for error
		widthMax       int
		signalChan     <-chan os.Signal
		pager          bool
		pagerLimit     int
		pagerCommand   string
		schema         string
	}

	runnerOptionsFn func(*runnerOption)
)

func NewRunner(opts ...runnerOptionsFn) (*Runner, error) {
	var err error
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

	r.client, err = nebula.NewNebulaClient(r.option.address, r.option.user, r.option.password,
		nebula.WithClientRequestTimeout(time.Duration(r.option.timeoutSec)*time.Second),
	)
	if err != nil {
		return nil, err
	}
	id, err := r.client.GetSessionId()
	if err != nil {
		return nil, err
	}
	r.sessionId = id
	if err := r.client.Ping(); err != nil {
		return nil, err
	}
	// set schema for playing data
	if r.option.schema != "" {
		if _, err := r.client.Execute(fmt.Sprintf("SESSION SET SCHEMA `%s`", r.option.schema)); err != nil {
			return nil, err
		}
	}
	var c cli.Cli
	if r.option.interactive {
		historyFile := filepath.Join(r.option.historyDir, ".nebula_history")
		c = cli.NewiCli(historyFile, r.option.user, r.option.enableGoPrompt)
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

func WithGoPrompt(enable bool) runnerOptionsFn {
	return func(o *runnerOption) {
		o.enableGoPrompt = enable
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
func WithSchema(schema string) runnerOptionsFn {
	return func(o *runnerOption) {
		o.schema = schema
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
// We treat one line as one query
// Add line break yourself as `SHOW \<CR>HOSTS`
func (r *Runner) loop() error {
	for {
		line, exit, err := r.cli.ReadLine()
		if err != nil {
			return err
		}
		if exit { // Ctrl+D
			r.printBoth("\n")
			return nil
		}
		if len(line) == 0 { // 1). The line input is empty, or 2). user presses ctrlC so the input is truncated
			continue
		}
		line = strings.TrimSpace(line)
		// record in file
		r.printFile(fmt.Sprintf("%s%s\n", r.cli.GetPrompt(), line))
		r.running = true
		exit, err = r.execute(line)
		r.running = false
		if err != nil {
			if r.option.failFast {
				return err
			}
			r.printBoth(fmt.Sprintf("Error: %s\n", err.Error()))
			continue
		}
		if exit { // :exit
			return nil
		}
	}
}

// execute one line
func (r *Runner) execute(line string) (exit bool, err error) {
	cmd, err := getCommand(r, line)
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
	isVertical := strings.HasSuffix(line, "\\G")
	if isVertical {
		line = strings.TrimSuffix(line, "\\G")
	}
	start := time.Now()
	resp, err := r.client.Execute(line)
	if err != nil {
		if r.option.failFast {
			return true, err
		}
		// do not exit, and continue to execute next statement
		r.printBoth(fmt.Sprintf("Error: %s\n", err.Error()))
	}
	duration := time.Since(start)
	if resp == nil {
		return false, nil
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

	return false, nil
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
