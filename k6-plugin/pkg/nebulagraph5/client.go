package nebulagraph5

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/vesoft-inc/k6-plugin/pkg/common"

	nebula "github.com/vesoft-inc/nebula-ng-tools/golang"
)

type (
	// GraphPool nebula connection pool
	GraphPool struct {
		DataCh            chan common.Data
		OutputCh          chan []string
		Version           string
		csvStrategy       csvReaderStrategy
		initialized       bool
		clients           []*GraphClient
		channelBufferSize int
		Hosts             []string
		mutex             sync.Mutex
		clientGetter      graphClientGetter
		csvReader         common.ICsvReader
		graphOption       *common.GraphOption
	}

	graphClientGetter func(endpoint, username, password string) (nebula.Client, error)

	// GraphClient a wrapper for nebula client, could read data from DataCh
	GraphClient struct {
		Session  nebula.Client
		Pool     *GraphPool
		DataCh   chan common.Data
		username string
		password string
	}

	// Response a wrapper for nebula resultSet
	Response struct {
		ResultSet    nebula.Result
		err          error
		ResponseTime int32
	}

	csvReaderStrategy int

	output struct {
		timeStamp    int64
		nGQL         string
		latency      int64
		responseTime int32
		isSucceed    bool
		rows         int32
		errorMsg     string
		firstRecord  string
	}
)

var _ common.IGraphClient = &GraphClient{}
var _ common.IGraphClientPool = &GraphPool{}

const (
	// AllInOne read csv sequentially
	AllInOne csvReaderStrategy = iota
	// Separate read csv concurrently
	Separate
)

func formatOutput(o *output) []string {
	return []string{
		strconv.FormatInt(o.timeStamp, 10),
		o.nGQL,
		strconv.Itoa(int(o.latency)),
		strconv.Itoa(int(o.responseTime)),
		strconv.FormatBool(o.isSucceed),
		strconv.Itoa(int(o.rows)),
		o.firstRecord,
		o.errorMsg,
	}
}

var outputHeader []string = []string{
	"timestamp",
	"nGQL",
	"latency",
	"responseTime",
	"isSucceed",
	"rows",
	"firstRecord",
	"errorMsg",
}

// NewNebulaGraph New for k6 initialization.
func NewNebulaGraph() *GraphPool {
	return &GraphPool{
		clientGetter: func(endpoint string, username, password string) (nebula.Client, error) {
			conn, err := nebula.NewNebulaClient(endpoint, username, password)
			if err != nil {
				return nil, err
			}
			if err := conn.Ping(); err != nil {
				return nil, fmt.Errorf("Failed to ping: %s", err.Error())
			}
			return conn, nil
		},
	}
}

func (gp *GraphPool) SetOption(option *common.GraphOption) error {
	if gp.graphOption != nil {
		return nil
	}
	gp.graphOption = common.MakeDefaultOption(option)
	if err := common.ValidateOption(gp.graphOption); err != nil {
		return err
	}
	bs, _ := json.Marshal(gp.graphOption)
	fmt.Printf("testing option: %s\n", bs)
	return nil
}

// Init initializes nebula pool with address and concurrent, by default the bufferSize is 20000
func (gp *GraphPool) Init() (common.IGraphClientPool, error) {
	gp.mutex.Lock()
	defer gp.mutex.Unlock()
	if gp.initialized {
		return gp, nil
	}
	if err := gp.validate(gp.graphOption.Address); err != nil {
		return nil, err
	}
	gp.Hosts = strings.Split(gp.graphOption.Address, ",")
	gp.clients = make([]*GraphClient, 0, gp.graphOption.MaxSize)
	if gp.graphOption.Output != "" {
		channelBufferSize := gp.graphOption.OutputChannelSize
		gp.OutputCh = make(chan []string, channelBufferSize)
		writer := common.NewCsvWriter(gp.graphOption.Output, ",", outputHeader, gp.OutputCh)
		if err := writer.WriteForever(); err != nil {
			return nil, err
		}
	}
	if gp.graphOption.CsvPath != "" {
		gp.csvReader = common.NewCsvReader(
			gp.graphOption.CsvPath,
			gp.graphOption.CsvDelimiter,
			gp.graphOption.CsvWithHeader,
			gp.graphOption.CsvDataLimit,
		)
		gp.DataCh = make(chan common.Data, gp.graphOption.CsvChannelSize)
		if err := gp.csvReader.ReadForever(gp.DataCh); err != nil {
			return nil, err
		}
	}
	gp.initialized = true
	return gp, nil
}

func (gp *GraphPool) validate(address string) error {
	addrs := strings.Split(address, ",")
	if len(addrs) == 0 {
		return fmt.Errorf("Invalid address: %s", address)
	}
	for _, addr := range addrs {
		hostAndPort := strings.Split(addr, ":")
		if len(hostAndPort) != 2 {
			return fmt.Errorf("Invalid address: %s", addr)
		}
	}
	return nil
}

// Close closes the nebula pool
func (gp *GraphPool) Close() error {
	gp.mutex.Lock()
	defer gp.mutex.Unlock()
	// gp.Log.Println("begin close the nebula pool")
	for _, s := range gp.clients {
		if s != nil {
			if s.Session != nil {
				s.Session.Close()
			}
		}
	}
	return nil
}

// GetSession gets the session from pool
func (gp *GraphPool) GetSession() (common.IGraphClient, error) {
	gp.mutex.Lock()
	defer gp.mutex.Unlock()
	index := len(gp.clients) % len(gp.Hosts)
	client, err := gp.clientGetter(
		gp.Hosts[index],
		gp.graphOption.Username,
		gp.graphOption.Password,
	)

	if err != nil {
		return nil, err
	}
	client.SetRequestTimeout(time.Duration(gp.graphOption.TimeoutUs))
	// client.SetGraph(gp.graphOption.Space)
	s := &GraphClient{Session: client, Pool: gp, DataCh: gp.DataCh}
	gp.clients = append(gp.clients, s)
	return s, nil
}

func (gc *GraphClient) Open() error {
	return nil
}
func (gc *GraphClient) Close() error {
	gc.Session.Close()
	return nil
}

// GetData get data from csv reader
func (gc *GraphClient) GetData() (common.Data, error) {
	if gc.DataCh != nil && len(gc.DataCh) != 0 {
		if d, ok := <-gc.DataCh; ok {
			return d, nil
		}
	}
	return nil, fmt.Errorf("no Data at all")
}

// Execute executes nebula query
func (gc *GraphClient) Execute(stmt string) (common.IGraphResponse, error) {
	var (
		isSucceed  bool = true
		errMessage string
		err        error
		resp       nebula.Result
		rows       int32
		latency    int64
	)
	stmt = common.ProcessStmt(stmt)
	start := time.Now()
	resp, err = gc.Session.Execute(stmt)

	if err != nil {
		isSucceed = false
		errMessage = err.Error()
		rows = 0
		latency = 0
		gc.Session.Close()
		sess, err := gc.Pool.GetSession()
		if err != nil {
			return nil, err
		}
		newClient := sess.(*GraphClient)
		*gc = *newClient
	} else {
		rows = int32(resp.RowSize())
		latency = resp.Latency()
	}

	responseTime := int32(time.Since(start) / 1000)
	// output

	if gc.Pool.OutputCh != nil {
		var fr []string
		if rows != 0 {
			// print the first row of the result
			row, err := resp.Next()
			if err != nil {
				return nil, err
			}

			for _, v := range row.Values() {
				fr = append(fr, v.String())
			}
		}
		o := &output{
			timeStamp:    start.Unix(),
			nGQL:         stmt,
			latency:      latency,
			responseTime: responseTime,
			isSucceed:    isSucceed,
			rows:         rows,
			errorMsg:     errMessage,
			firstRecord:  strings.Join(fr, "|"),
		}
		select {
		case gc.Pool.OutputCh <- formatOutput(o):
		// abandon if the output chan is full.
		default:
			fmt.Printf("output channel is full, abandon the output: %v\n", o)
		}

	}
	return &Response{ResultSet: resp, ResponseTime: responseTime, err: err}, nil
}

// GetResponseTime GetResponseTime
func (r *Response) GetResponseTime() int32 {
	return r.ResponseTime
}

// IsSucceed IsSucceed
func (r *Response) IsSucceed() bool {
	if r.err != nil {
		return false
	}

	return true
}

func (r *Response) GetLatency() int64 {
	if r.ResultSet != nil {
		return r.ResultSet.Latency()
	}
	return 0
}

// GetRowSize GetRowSize
func (r *Response) GetRowSize() int32 {
	if r.ResultSet != nil {
		return int32(r.ResultSet.RowSize())
	}
	return 0
}
