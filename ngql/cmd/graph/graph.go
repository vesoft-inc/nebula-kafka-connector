// Copyright (c) 2023 vesoft inc. All rights reserved.
package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/vesoft-inc/nebula-ng-tools/ngql/pkg/runner"
)

// root flags variables
var (
	host           string
	port           int
	username       string
	password       string
	timeout        int
	script         string
	file           string
	version        bool
	widthMax       int
	enableTLS      bool
	ca             string
	cert           string
	key            string
	peerNameVerify bool
	peerName       string
	goPrompt       bool
)

type ParameterMap map[string]interface{}

var parameterMap ParameterMap

func validateFlags() {
	if port == -1 {
		log.Fatalf("Error: argument port is missing!")
	}
	if len(username) == 0 {
		log.Fatalf("Error: argument username is empty!")
	}

	if widthMax < 0 || (widthMax > 0 && widthMax <= 3) {
		log.Fatalf("Error: argument width_max should be equal to 0 or greater than 3")
	}

	if enableTLS {
		if ca == "" {
			log.Fatalf("Error: --ca must be specified if --enable-tls is true")
		}
		if (cert != "" && key == "") || (cert == "" && key != "") {
			log.Fatalf("Error: --cert and --key must be specified in pairs")
		}
	}
}

func handleGraphCmd() error {
	parameterMap = make(ParameterMap)

	// if flag.NFlag() == 1 && version {
	// 	fmt.Printf("ngql version Git: %s, Build Time: %s\n", gitCommit, buildDate)
	// 	return nil
	// }

	// Check if flags are valid
	validateFlags()

	interactive := script == "" && file == ""
	var fastFail bool
	var rc io.ReadCloser
	output := true
	if !interactive {
		fastFail = true
		if file != "" {
			fd, err := os.Open(file)
			if err != nil {
				log.Fatalf("Open file %s failed, %s", file, err.Error())
			}
			rc = fd
			// If file is provided, we should not output the result
			output = false
		} else {
			rc = ioutil.NopCloser(strings.NewReader(script))
		}
	}
	address := fmt.Sprintf("%s:%d", host, port)
	historyHome := os.Getenv("HOME")
	if historyHome == "" {
		ex, err := os.Executable()
		if err != nil {
			log.Fatalf("Get executable failed: %s", err.Error())
		}
		historyHome = filepath.Dir(ex) // Set to executable folder
	}
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	runner, err := runner.NewRunner(
		runner.WithInteractive(interactive),
		runner.WithNebula(address, username, password),
		runner.WithTimeoutSec(timeout),
		runner.WithHistoryDir(historyHome),
		runner.WithReadCloser(rc),
		runner.WithGoPrompt(goPrompt),
		runner.WithOutput(output),
		runner.WithWidthMax(widthMax),
		runner.WithSignalChan(c),
		runner.WithFailFast(fastFail),
		runner.WithTLS(enableTLS, ca, cert, key, peerNameVerify, peerName),
	)
	if err != nil {
		return err
	}
	defer runner.Close()

	if err := runner.Run(); err != nil {
		return err
	}
	if !interactive {
		fmt.Fprintln(os.Stdout, "Execute successfully!")
	}
	return nil

}

var rootCmd = &cobra.Command{
	Use:   "ngql",
	Short: "Run ngql to connect to nebula-graphd by default.",
	Long:  `Use ngql --host [host] --port [port] -u [user] -p [password] to connect to nebula-graphd.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if password == "" {
			var err error
			pw := promptui.Prompt{
				Label:       "Password",
				AllowEdit:   true,
				Mask:        rune(' '),
				HideEntered: true,
			}
			password, err = pw.Run()
			if err != nil {
				return err
			}
		}
		return handleGraphCmd()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&host, "host", "H", "127.0.0.1", "Nebula Graph host")
	rootCmd.Flags().IntVarP(&port, "port", "P", 9669, "The Nebula Graph port")
	rootCmd.Flags().StringVarP(&username, "user", "u", "root", "The Nebula Graph login user name")
	rootCmd.Flags().StringVarP(&password, "password", "p", "", "The Nebula Graph login password")
	rootCmd.Flags().IntVarP(&timeout, "timeout", "t", 0, "The Nebula Graph client connection timeout in seconds, 0 means never timeout")
	rootCmd.Flags().StringVarP(&script, "eval", "e", "", "The GQL directly")
	rootCmd.Flags().StringVarP(&file, "file", "f", "", "The GQL script file name")
	rootCmd.Flags().IntVarP(&widthMax, "width-max", "", 100, "The max width of the column of the execution plan")
	rootCmd.Flags().BoolVarP(&enableTLS, "enable-tls", "", false, "Enable TLS")
	rootCmd.Flags().StringVarP(&ca, "ca", "", "", "Certificate(s) of trusted CA, in PEM format")
	rootCmd.Flags().StringVarP(&cert, "cert", "", "", "Client certificate, in PEM format")
	rootCmd.Flags().StringVarP(&key, "key", "", "", "Client private key, in PEM format")
	rootCmd.Flags().BoolVarP(&peerNameVerify, "peer-name-verify", "", false, "Enable peer name verification")
	rootCmd.Flags().StringVarP(&peerName, "peer-name", "", "", "Peer name to override the default, i.e. domain name")
}

func main() {
	rootCmd.Execute()
}
