package runner

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	nebula "github.com/vesoft-inc/nebula-ng-tools/golang"
	"github.com/vesoft-inc/nebula-ng-tools/ngql/pkg/cli"
	"github.com/vesoft-inc/nebula-ng-tools/ngql/pkg/printer"
)

type (
	Runner struct {
		option       *runnerOption
		stdout       io.Writer
		file         io.WriteCloser
		client       nebula.Client
		cli          cli.Cli
		printer      printer.Printer
		combinOutput io.Writer
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
	}

	runnerOptionsFn func(*runnerOption)
)

func NewRunner(opts ...runnerOptionsFn) (*Runner, error) {
	var err error
	r := &Runner{
		stdout:  os.Stdout,
		option:  &runnerOption{},
		printer: printer.NewPrinter("default", 100),
	}
	r.combinOutput = &combinOutput{
		runner: r,
	}
	for _, opt := range opts {
		opt(r.option)
	}
	if !r.option.enableOutput {
		r.stdout = io.Discard
	}

	r.client, err = nebula.NewNebulaClient(r.option.address, r.option.user, r.option.password)
	if err != nil {
		return nil, err
	}
	r.client.SetRequestTimeout(time.Duration(r.option.timeoutSec) * time.Second)
	if err := r.client.Ping(); err != nil {
		return nil, err
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

func (r *Runner) welcome() {
	if !r.option.interactive {
		return
	}
	r.printStdout("Welcome to Nebula Graph!\n")
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
		exit, err = r.execute(line)
		if err != nil {
			r.printBoth(fmt.Sprintf("Error: %s\n", err.Error()))
			if r.option.failFast {
				return err
			}
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
		if cmd.cmd == commandExit {
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
	if resp == nil {
		return false, nil
	}

	if isVertical {
		r.printer.PrintResultVertical(r.combinOutput, resp)
	} else {
		r.printer.PrintResult(r.combinOutput, resp)
	}
	duration := time.Since(start)
	if resp.RowSize() > 0 {
		numRows := resp.RowSize()
		r.printBoth(fmt.Sprintf("Got %d rows (time spent %v/%v)\n\n",
			numRows, time.Duration(resp.Latency()*1000), duration),
		)
	} else {
		r.printBoth(fmt.Sprintf("Execution succeeded (time spent %v/%v)\n\n",
			time.Duration(resp.Latency()*1000), duration),
		)
	}
	if resp.Summary() != nil {
		r.printer.PrintPlanDesc(r.combinOutput, resp.Summary())
	}
	r.printBoth(fmt.Sprintf("\n%s\n\n", time.Now().In(time.Local).Format(time.RFC1123)))

	return false, nil
}

func (r *Runner) Run() error {
	if r.option.interactive {
		r.welcome()
		defer r.bye()
	}
	if err := r.loop(); err != nil {
		r.printBoth(fmt.Sprintf("Error: %s\n", err))
		return err
	}
	return nil
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

func (o *combinOutput) Write(p []byte) (n int, err error) {
	n, err = o.runner.stdout.Write(p)
	if o.runner.file != nil {
		o.runner.file.Write(p)
	}
	return
}
