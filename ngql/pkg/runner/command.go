package runner

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/vesoft-inc/nebula-ng-tools/ngql/pkg/data"
	"github.com/vesoft-inc/nebula-ng-tools/ngql/pkg/printer"
)

type command struct {
	runner   *Runner  // used for some command to change runner context
	args     []string // e.g. []string{"1"}
	fn       commandFn
	cmd      string
	desc     string
	usage    string
	alias    string
	showHelp bool
}

const (
	commandSleep   = "sleep"
	commandPlay    = "play"
	commandRepeat  = "repeat"
	commandFormat  = "format"
	commandExit    = "exit"
	commandQuit    = "quit"
	commandHelp    = "help"
	commnadTee     = "tee"
	commnadNoTee   = "notee"
	commandPager   = "pager"
	commandNoPager = "nopager"
)

type commandFn func(r *Runner, args []string) error
type commandOptionFn func(c *command)

// name for commands, used for help
var commandNames []string
var commands map[string]*command

func (c *command) setArgs(args []string) {
	c.args = args
}
func (c *command) setRunner(r *Runner) {
	c.runner = r
}

func (c *command) execute() error {
	return c.fn(c.runner, c.args)
}

func withShowHelp(showHelp bool) commandOptionFn {
	return func(c *command) {
		c.showHelp = showHelp
	}
}

func registerCommand(cmd string, usage, alias, desc string, fn commandFn, options ...commandOptionFn) {
	if commands == nil {
		commands = make(map[string]*command)
	}
	if commandNames == nil {
		commandNames = make([]string, 0)
	}
	commandNames = append(commandNames, cmd)
	c := &command{
		cmd:      cmd,
		alias:    alias,
		desc:     desc,
		usage:    usage,
		fn:       fn,
		showHelp: true, // show help by default
	}
	for _, option := range options {
		option(c)
	}

	commands[cmd] = c
	if alias != "" {
		a := strings.Trim(alias, ":")
		commands[a] = c
	}
}

func getCommand(r *Runner, line string) (*command, error) {
	if !strings.HasPrefix(line, ":") {
		return nil, nil
	}
	l := strings.Trim(line, ":")
	l = strings.TrimSpace(l)
	words := strings.Fields(l)
	if len(words) == 0 {
		return nil, nil
	}
	cmd := words[0]
	cmd = strings.ToLower(cmd)
	args := words[1:]
	c, ok := commands[cmd]
	if !ok {
		return nil, fmt.Errorf("Command not found: %s", cmd)
	}
	c.setArgs(args)
	c.setRunner(r)
	return c, nil

}

func sleepFn(r *Runner, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("Sleep command needs an argument")
	}
	i, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}
	time.Sleep(time.Duration(i) * time.Second)
	return nil
}
func helpFn(r *Runner, args []string) error {
	writer := table.NewWriter()
	writer.Style().Format.Header = text.FormatDefault
	writer.Style().Options.SeparateRows = false
	writer.ResetHeaders()
	writer.ResetRows()
	writer.AppendHeader(table.Row{"Command", "Alias", "Usage", "Description"})
	for _, name := range commandNames {
		c := commands[name]
		if !c.showHelp {
			continue
		}
		writer.AppendRow(table.Row{c.cmd, c.alias, c.usage, c.desc})
	}
	writer.Render()
	fmt.Fprintf(r.stdout, writer.Render())
	fmt.Fprintf(r.stdout, "\n")
	return nil
}

func playFn(r *Runner, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("Play command needs an argument")
	}
	playName := args[0]
	// TODO maybe we could download the file from oss.
	// and then play the file.
	entries, err := data.Play.ReadDir(".")
	if err != nil {
		return err
	}
	var found = false
	var filename string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		filename = entry.Name()
		if !strings.HasSuffix(filename, ".ngql") {
			continue
		}
		datasetName := strings.TrimSuffix(filename, ".ngql")

		if datasetName == playName {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("Play file not found: %s", playName)
	}

	bs, err := data.Play.ReadFile(filename)
	if err != nil {
		return err
	}
	//get current session schema
	resp, err := r.client.Execute("SHOW CURRENT_SESSION")
	if err != nil {
		return err
	}
	var scehma string
	for resp.HasNext() {
		row, err := resp.Next()
		if err != nil {
			return err
		}
		v, err := row.GetValueByName("home_schema_path")
		if err != nil {
			return err
		}
		ss, err := v.AsString()
		if err != nil {
			return err
		}
		scehma = string(ss)
		break
	}

	subRunner, err := NewRunner(
		WithInteractive(false),
		WithReadCloser(io.NopCloser(bytes.NewBuffer(bs))),
		WithOutput(false),
		WithFailFast(true),
		WithNebula(r.option.address, r.option.user, r.option.password),
	)
	if err != nil {
		return err
	}
	// set playing schema
	_, err = r.client.Execute(fmt.Sprintf("SESSION SET SCHEMA \"%s\"", scehma))
	if err != nil {
		return err
	}
	fmt.Fprintf(r.stdout, "Playing dataset: %s...\n", playName)
	if err := subRunner.Run(); err != nil {
		return err
	} else {
		fmt.Fprintf(r.stdout, "Play dataset: %s done.\n", playName)
		return nil
	}
}

func formatFn(r *Runner, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("Format command needs an argument")
	}
	r.printBoth(fmt.Sprintf("Set format to %s\n", args[0]))
	r.printer = printer.NewPrinter(args[0], r.option.widthMax)
	return nil
}

func teeFn(r *Runner, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("Tee command needs an argument")
	}
	var overwrite bool
	var filename string
	if len(args) == 1 {
		filename = args[0]
	} else if len(args) == 2 && args[0] == "-o" {
		overwrite = true
		filename = args[1]
	} else {
		return fmt.Errorf("Invalid tee command, usage: tee [-o] filename")
	}
	var (
		file *os.File
		err  error
	)
	if overwrite {
		file, err = os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	} else {
		file, err = os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	}
	if err != nil {
		return err
	}
	r.file = file
	return nil
}

func noteeFn(r *Runner, args []string) error {
	if r.file != nil {
		r.file.Close()
		r.file = nil
	}
	return nil
}

func pagerFn(r *Runner, args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("Pager command needs two arguments, e.g. :pager less 200")
	}
	i, err := strconv.Atoi(args[1])
	if err != nil {
		return err
	}
	//check if the command is valid
	c := args[0]
	cmd := exec.Command(c)
	if err := cmd.Run(); err != nil {
		return err
	}

	r.option.pagerLimit = i
	r.option.pagerCommand = c
	r.option.pager = true
	r.printStdout(fmt.Sprintf("Pager set to %s with row limit %d\n", args[0], i))
	return nil
}

func noPagerFn(r *Runner, args []string) error {
	r.option.pager = false
	r.printStdout("Pager disabled.\n")
	return nil
}

func init() {
	registerCommand(commandHelp, ":help", ":h", "Show this help.", helpFn)
	registerCommand(commandSleep, ":sleep 5", "", "Sleep N seconds.", sleepFn)
	registerCommand(commandPlay, ":play movie", "", "Playing the dateset", playFn)
	registerCommand(commandFormat, ":format default", "", "Change the format for value. (default, tck)", formatFn, withShowHelp(false))
	registerCommand(commnadTee, ":tee [-o] <filename>", "", "Append all results to an output file (overwrite using -o).", teeFn)
	registerCommand(commnadNoTee, ":notee", "", "Stop writing to the output file.", noteeFn)
	registerCommand(commandPager, ":pager <commnad> <row_limit>", "", "Set pager for result, default: \":pager less 200\"", pagerFn)
	registerCommand(commandNoPager, ":nopager", "", "No pager", noPagerFn)
	registerCommand(commandExit, ":exit", ":e", "Exit.", func(r *Runner, args []string) error {
		return nil
	})
	registerCommand(commandQuit, ":quit", ":q", "Quit.", func(r *Runner, args []string) error {
		return nil
	})

}
