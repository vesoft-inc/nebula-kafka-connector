// Copyright (c) 2023 vesoft inc. All rights reserved.
package cmd

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/vesoft-inc/nebula-ng-tools/console/box"
	"github.com/vesoft-inc/nebula-ng-tools/console/cli"
	"github.com/vesoft-inc/nebula-ng-tools/console/printer"
	nebulago "github.com/vesoft-inc/nebula-ng-tools/golang"
)

// root flags variables
var (
	address               string
	port                  int
	username              string
	password              string
	timeout               int
	script                string
	file                  string
	version               bool
	widthMax              int
	enableSsl             bool
	sslRootCAPath         string
	sslCertPath           string
	sslPrivateKeyPath     string
	sslInsecureSkipVerify bool
	goPrompt              bool
)

// Console side commands
const (
	Unknown   = -1
	Quit      = 0
	PlayData  = 1
	Sleep     = 2
	ExportCsv = 3
	ExportDot = 4
	Repeat    = 5
	Param     = 6
	Params    = 7
)

type ParameterMap map[string]interface{}

var parameterMap ParameterMap

var dataSetPrinter = printer.NewDataSetPrinter()

var planDescPrinter = printer.NewPlanDescPrinter()

/*
	Every statement will be repeatedly executed `g_repeats` times,

in order to get the total and average execution time of the statement")
*/
var g_repeats = 1

func welcome(interactive bool) {
	if !interactive {
		return
	}
	fmt.Println()
	fmt.Printf("Welcome to Nebula Graph!\n")
	fmt.Println()
}

func bye(username string, interactive bool) {
	fmt.Println()
	fmt.Printf("Bye %s!\n", username)
	fmt.Println(time.Now().In(time.Local).Format(time.RFC1123))
	fmt.Println()
}

func printConsoleResp(msg string) {
	fmt.Println(msg)
	fmt.Println()
	fmt.Println(time.Now().In(time.Local).Format(time.RFC1123))
	fmt.Println()
}

func playData(data string) error {
	boxfilePath := "/" + data + ".ngql"
	posixfilePath := "./data/" + data + ".ngql"
	var c cli.Cli
	// First find it in directory ./data/. If not found, then find it in the embedded box
	if fd, err := os.Open(posixfilePath); err == nil {
		c = cli.NewnCli(fd, false, "", func() { fd.Close() })
	} else if box.Has(boxfilePath) {
		fileStr := string(box.Get(boxfilePath))
		c = cli.NewnCli(strings.NewReader(fileStr), false, "", nil)
	} else {
		return fmt.Errorf("file %s.ngql not existed in embed box and file directory ./data/ ", data)
	}

	c.PlayingData(true)
	defer c.PlayingData(false)
	fmt.Printf("Start loading dataset %s...\n", data)
	err := loop(c)
	if err != nil {
		return err
	}
	respErr := c.GetRespError()
	if respErr != "" {
		return fmt.Errorf(respErr)
	}
	return nil
}

func defineParams(args string) {
	argsRewritten := strings.Replace(args, "'", "\"", -1)
	reg := regexp.MustCompile(`^\s*:param\s+(\S+)\s*=>(.*)$`)
	if reg == nil {
		fmt.Println("invalid regular expression")
		return
	}
	matchResult := reg.FindAllStringSubmatch(argsRewritten, -1)
	if len(matchResult) != 1 || len(matchResult[0]) != 3 {
		fmt.Println("Wrong local command format", matchResult)
		return
	}
	/*
	 * :param p1=> -> [":param p1=>",":p1",""]
	 * :param p2=>3 -> [":param p2=>3",":p2","3"]
	 */
	paramKey := matchResult[0][1]
	paramValue := matchResult[0][2]
	if len(paramValue) == 0 {
		delete(parameterMap, paramKey)
	} else {
		paramsWithGoType := make(ParameterMap)
		param := "{\"" + paramKey + "\"" + ":" + paramValue + "}"
		err := json.Unmarshal([]byte(param), &paramsWithGoType)
		if err != nil {
			fmt.Println("Error: parameter parsing failed")
			return
		}
		for k, v := range paramsWithGoType {
			parameterMap[k] = v
		}
	}
}

func ListParams(args string) {
	reg := regexp.MustCompile(`^\s*:params\s*(\S*)\s*$`)
	if reg == nil {
		fmt.Println("invalid regular expression")
		return
	}
	matchResult := reg.FindAllStringSubmatch(args, -1)
	if len(matchResult) != 1 {
		fmt.Println("Wrong local command format", matchResult)
		return
	}
	res := matchResult[0]
	/*
	 * :params -> [":params",""]
	 * :params p1 -> ["params","p1"]
	 */
	if len(res) != 2 {
		return
	} else {
		paramKey := matchResult[0][1]
		if len(paramKey) == 0 {
			for k, v := range parameterMap {
				fmt.Println(k, " => ", v)
			}
		} else {
			if paramValue, ok := parameterMap[paramKey]; ok {
				fmt.Println(paramKey, " => ", paramValue)
			} else {
				fmt.Println("Unknown parameter: ", paramKey)
			}
		}
	}
}

// Console side cmd will not be sent to server
func isConsoleCmd(cmd string) (isLocal bool, localCmd int, args []string) {
	isLocal = false
	localCmd = Unknown
	// Currently, command "exit" and  "quit" can also exit the console
	if cmd == "exit" || cmd == "quit" {
		isLocal = true
		localCmd = Quit
		return
	}

	plain := strings.TrimSpace(cmd)
	if len(plain) < 1 || plain[0] != ':' {
		return
	}

	isLocal = true
	if plain[len(plain)-1] == ';' {
		plain = plain[:len(plain)-1]
	}
	words := strings.Fields(plain[1:])
	localCmdName := words[0]
	switch strings.ToLower(localCmdName) {
	case "exit", "quit":
		{
			localCmd = Quit
		}
	case "sleep":
		{
			localCmd = Sleep
			args = []string{words[1]}
		}
	case "play":
		{
			localCmd = PlayData
			args = []string{words[1]}
		}
	case "repeat":
		{
			localCmd = Repeat
			args = []string{words[1]}
		}
	case "csv":
		{
			localCmd = ExportCsv
			args = []string{words[1]}
		}
	case "dot":
		{
			localCmd = ExportDot
			args = []string{words[1]}
		}
	case "param":
		{
			localCmd = Param
			args = []string{plain}
		}
	case "params":
		{
			localCmd = Params
			args = []string{plain}
		}
	}
	return
}

func executeConsoleCmd(c cli.Cli, cmd int, args []string) {
	switch cmd {
	case ExportCsv:
		dataSetPrinter.ExportCsv(args[0])
	case ExportDot:
		planDescPrinter.ExportDot(args[0])
	case PlayData:
		err := playData(args[0])
		if err != nil {
			printConsoleResp("Error: load dataset failed, " + err.Error())
		} else {
			printConsoleResp("Load dataset succeeded!")
			// c.SetSpace(newSpace)
		}
	case Sleep:
		i, err := strconv.Atoi(args[0])
		if err != nil {
			printConsoleResp("Error: invalid integer, " + err.Error())
		}
		time.Sleep(time.Duration(i) * time.Second)
	case Repeat:
		i, err := strconv.Atoi(args[0])
		if err != nil {
			printConsoleResp("Error: invalid integer, " + err.Error())
		} else if i < 1 {
			printConsoleResp("Error: invald integer, repeats should be greater than 1")
		}
		g_repeats = i
	case Param:
		if len(args) != 1 {
			return
		}
		defineParams(args[0])
	case Params:
		if len(args) != 1 {
			return
		}
		ListParams(args[0])
	default:
		printConsoleResp("Error: this local command not exists!")
	}
}

// TODO(Aiee) We don't have a complete gql status yet
func printResultSet(res *nebulago.ResultSet, startTime time.Time) (duration time.Duration) {
	if !res.IsSucceed() {
		// fmt.Printf("[ERROR (%d)]: %s", res.GetErrorCode(), res.GetErrorMsg())
		fmt.Printf("[ERROR]: %s", res.GetStatus())
		fmt.Println()
		fmt.Println()
		return
	}
	// Show table
	if res.IsSetData() {
		dataSetPrinter.PrintDataSet(res)
		numRows := res.GetRowSize()
		duration = time.Since(startTime)
		if numRows > 0 {
			fmt.Printf("Got %d rows (time spent %v/%v)\n", numRows, time.Duration(res.GetLatency()*1000), duration)
		} else {
			fmt.Printf("Empty set (time spent %v/%v)\n", time.Duration(res.GetLatency()*1000), duration)
		}
	} else {
		duration = time.Since(startTime)
		fmt.Printf("Execution succeeded (time spent %v/%v)\n", time.Duration(res.GetLatency()*1000), duration)
	}

	// if res.IsPartialSucceed() {
	// 	fmt.Println()
	// 	fmt.Printf("[WARNING]: Got partial result.")
	// }

	// if res.IsSetComment() {
	// 	fmt.Println()
	// 	fmt.Printf("[WARNING]: %s", res.GetComment())
	// }

	if res.IsSetPlanDesc() {
		fmt.Println()
		planDescPrinter.PrintPlanDesc(res)
	}
	fmt.Println()

	return
}

// Loop the request util fatal or timeout
// We treat one line as one query
// Add line break yourself as `SHOW \<CR>HOSTS`
func loop(c cli.Cli) error {
	for {
		line, exit, err := c.ReadLine()
		if err != nil {
			return err
		}
		if exit { // Ctrl+D
			fmt.Println()
			return nil
		}
		if len(line) == 0 { // 1). The line input is empty, or 2). user presses ctrlC so the input is truncated
			continue
		}
		// Console side command
		if isLocal, cmd, args := isConsoleCmd(line); isLocal {
			if cmd == Quit {
				return nil
			}
			executeConsoleCmd(c, cmd, args)
			continue
		}
		// Server side command
		var t1 int64 = 0
		var t2 int64 = 0
		for i := 0; i < g_repeats; i++ {
			start := time.Now()
			res, err := session.Execute(line)
			if err != nil {
				return err
			}

			if c.Output() {
				duration := printResultSet(res, start)
				t2 += int64(duration / 1000)
				fmt.Println(time.Now().In(time.Local).Format(time.RFC1123))
				fmt.Println()
			}
		}
		if g_repeats > 1 {
			fmt.Printf(
				"Executed %v times, (total time spent %d/%d us), (average time spent %d/%d us)\n",
				g_repeats,
				t1,
				t2,
				t1/int64(g_repeats),
				t2/int64(g_repeats),
			)
			fmt.Println()
		}
		g_repeats = 1
	}
}

func openAndReadFile(path string) ([]byte, error) {
	// open file
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("unable to open file %s: %s", path, err)
	}
	// read file
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("unable to ReadAll of file %s: %s", path, err)
	}
	return b, nil
}

func genSslConfig(rootCAPath, certPath, privateKeyPath string) (*tls.Config, error) {
	rootCA, err := openAndReadFile(rootCAPath)
	if err != nil {
		return nil, err
	}
	cert, err := openAndReadFile(certPath)
	if err != nil {
		return nil, err
	}
	privateKey, err := openAndReadFile(privateKeyPath)
	if err != nil {
		return nil, err
	}

	// generate the client certificate
	clientCert, err := tls.X509KeyPair(cert, privateKey)
	if err != nil {
		return nil, err
	}

	// parse root CA pem and add into CA pool
	rootCAPool := x509.NewCertPool()
	ok := rootCAPool.AppendCertsFromPEM(rootCA)
	if !ok {
		return nil, fmt.Errorf("fail to append supplied cert into tls.Config, please make sure it is a valid certificate")
	}

	// set tls config
	// InsecureSkipVerify is set to true for test purpose ONLY. DO NOT use it in production.
	return &tls.Config{
		Certificates:       []tls.Certificate{clientCert},
		RootCAs:            rootCAPool,
		InsecureSkipVerify: sslInsecureSkipVerify,
	}, nil
}

// Nebula Console version related
var (
	gitCommit string
	buildDate string
)

func validateFlags() {
	if port == -1 {
		log.Panicf("Error: argument port is missed!")
	}
	if len(username) == 0 {
		log.Panicf("Error: argument username is empty!")
	}
	if len(password) == 0 {
		log.Panicf("Error: argument password is empty!")
	}

	if widthMax < 0 || (widthMax > 0 && widthMax <= 3) {
		log.Panicf("Error: argument width_max should be equal to 0 or greater than 3")
	}

	if enableSsl {
		if sslRootCAPath == "" {
			log.Panicf("Error: argument ssl_root_ca_path should be specified when enable_ssl is true")
		}
		if sslCertPath == "" {
			log.Panicf("Error: argument ssl_cert_path should be specified when enable_ssl is true")
		}
		if sslPrivateKeyPath == "" {
			log.Panicf("Error: argument ssl_private_key_path should be specified when enable_ssl is true")
		}
	}
}

var session *nebulago.Session

func handleGraphCmd() {
	parameterMap = make(ParameterMap)

	if flag.NFlag() == 1 && version {
		fmt.Printf("nebula-console version Git: %s, Build Time: %s\n", gitCommit, buildDate)
		return
	}

	// Check if flags are valid
	validateFlags()
	planDescPrinter.WidthMax = widthMax

	interactive := script == "" && file == ""

	historyHome := os.Getenv("HOME")
	if historyHome == "" {
		ex, err := os.Executable()
		if err != nil {
			log.Panicf("Get executable failed: %s", err.Error())
		}
		historyHome = filepath.Dir(ex) // Set to executable folder
	}

	hostAddr := nebulago.HostAddress{Host: address, Port: port}
	// hostList := []nebulago.HostAddress{hostAddress}
	// poolConfig := nebulago.PoolConfig{
	// 	TimeOut:         time.Duration(*timeout) * time.Millisecond,
	// 	IdleTime:        0 * time.Millisecond,
	// 	MaxConnPoolSize: 2,
	// 	MinConnPoolSize: 1,
	// }
	var err error
	// if *enableSsl {
	// 	sslConfig, err2 := genSslConfig(*sslRootCAPath, *sslCertPath, *sslPrivateKeyPath)
	// 	if err2 != nil {
	// 		log.Panicf(fmt.Sprintf("Fail to generate the ssl config, ssl_root_ca_path: %s, ssl_cert_path: %s, ssl_private_key_path: %s, %s", *sslRootCAPath, *sslCertPath, *sslPrivateKeyPath, err2.Error()))
	// 	}
	// 	pool, err = nebulago.NewSslConnectionPool(hostList, poolConfig, sslConfig, nebulago.DefaultLogger{})
	// } else {
	// 	pool, err = nebulago.NewConnectionPool(hostList, poolConfig, nebulago.DefaultLogger{})
	// }
	// if err != nil {
	// 	log.Panicf(fmt.Sprintf("Fail to initialize the connection pool, host: %s, port: %d, %s", *address, *port, err.Error()))
	// }
	// defer pool.Close()

	// Build connection
	connection := nebulago.NewConnection(hostAddr)
	err = connection.Open(hostAddr, time.Duration(timeout)*time.Millisecond, nil)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Authenticate to get the identifier
	authResp, err := connection.Authenticate(username, password)
	if err != nil {
		log.Fatal(err.Error())
	}
	if string(authResp.GetGqlStatus().Status) != "SUCCESS" {
		log.Fatal(fmt.Sprintf("authentication failed, error: %s", string(authResp.GetGqlStatus().Status)))
	}
	log.Println(fmt.Sprintf("Authentication with Identifier: %d succeed", authResp.GetIdentifier()))

	// Build session
	session = nebulago.NewSession(authResp.GetIdentifier(), connection, nebulago.DefaultLogger{})
	// if err != nil {
	// 	log.Panicf("Fail to create a new session from connection pool, %s", err.Error())
	// }
	defer session.Release()

	welcome(interactive)
	defer bye(username, interactive)

	var c cli.Cli = nil
	// Loop the request
	if interactive {
		historyFile := path.Join(historyHome, ".nebula_history")
		c = cli.NewiCli(historyFile, username, goPrompt)
	} else if script != "" {
		c = cli.NewnCli(strings.NewReader(script), true, username, nil)
	} else if file != "" {
		fd, err := os.Open(file)
		if err != nil {
			log.Panicf("Open file %s failed, %s", file, err.Error())
		}
		c = cli.NewnCli(fd, true, username, func() { fd.Close() })
	}

	if c == nil {
		return
	}

	defer c.Close()

	err = loop(c)

	if err != nil {
		log.Panicf("Loop error, %s", err.Error())
	}
}

var rootCmd = &cobra.Command{
	Use:   "nebula-console",
	Short: "Run nebula-console to connect to nebula-graphd by default.",
	Long:  `Use nebula-console --addr [addr] --port [port] -u [user] -p [password] to connect to nebula-graphd.`,
	Run: func(cmd *cobra.Command, args []string) {
		handleGraphCmd()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&address, "addr", "", "127.0.0.1", "Nebula Graph address")
	rootCmd.Flags().IntVarP(&port, "port", "", -1, "The Nebula Graph port")
	rootCmd.Flags().StringVarP(&username, "user", "u", "", "The Nebula Graph login user name")
	rootCmd.Flags().StringVarP(&password, "password", "p", "", "The Nebula Graph login password")
	rootCmd.Flags().IntVarP(&timeout, "timeout", "t", 0, "The Nebula Graph client connection timeout in seconds, 0 means never timeout")
	rootCmd.Flags().StringVarP(&script, "eval", "e", "", "The nGQL directly")
	rootCmd.Flags().StringVarP(&file, "file", "f", "", "The nGQL script file name")
	rootCmd.Flags().IntVarP(&widthMax, "width_max", "", 100, "The max width of the column of the execution plan")
}
