package nebulagraph5

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/vesoft-inc/nebula-ng-tools/k6-plugin/pkg/common"

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
		clients           []*nebula.Session
		channelBufferSize int
		hosts             []string
		mutex             sync.Mutex
		clientGetter      graphClientGetter
		csvReader         common.ICsvReader
	}

	graphClientGetter func(endpoint, username, password string) (*nebula.Session, error)

	// GraphClient a wrapper for nebula client, could read data from DataCh
	GraphClient struct {
		Client   *nebula.Session
		Pool     *GraphPool
		DataCh   chan common.Data
		username string
		password string
	}

	// Response a wrapper for nebula resultSet
	Response struct {
		*nebula.ResultSet
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
		clientGetter: func(endpoint string, username, password string) (*nebula.Session, error) {
			if len(strings.Split(endpoint, ":")) != 2 {
				return nil, fmt.Errorf("Invalid address: %s", endpoint)
			}
			host, p := strings.Split(endpoint, ":")[0], strings.Split(endpoint, ":")[1]
			port, err := strconv.Atoi(p)
			if err != nil {
				return nil, err
			}
			hostAddr := nebula.HostAddress{Host: host, Port: port}
			conn := nebula.NewConnection(hostAddr)
			if err := conn.Open(hostAddr, 3*time.Second, nil); err != nil {
				return nil, fmt.Errorf("Failed to open connection: %s", err.Error())
			}
			// Authenticate to get the identifier
			authResp, err := conn.Authenticate(username, password)
			if err != nil {
				return nil, fmt.Errorf("Failed to authenticate: %s", err.Error())
			}
			if string(authResp.GetGqlStatus().Status) != "SUCCESS" {
				return nil, fmt.Errorf("authentication failed, error: %s", string(authResp.GetGqlStatus().Status))
			}
			session := nebula.NewSession(authResp.GetIdentifier(), conn, &nebula.DefaultLogger{})
			return session, nil
		},
	}
}

// Init initializes nebula pool with address and concurrent, by default the bufferSize is 20000
func (gp *GraphPool) Init(address string, concurrent int) (common.IGraphClientPool, error) {
	return gp.InitWithSize(address, concurrent, 20000)
}

// InitWithSize initializes nebula pool with channel buffer size
func (gp *GraphPool) InitWithSize(address string, concurrent int, chanSize int) (common.IGraphClientPool, error) {
	gp.mutex.Lock()
	defer gp.mutex.Unlock()
	if gp.initialized {
		return gp, nil
	}

	err := gp.initAndVerifyPool(address, concurrent, chanSize)
	if err != nil {
		return nil, err
	}
	gp.DataCh = make(chan common.Data, chanSize)
	gp.initialized = true

	return gp, nil
}

func (gp *GraphPool) initAndVerifyPool(address string, concurrent int, chanSize int) error {
	addrs := strings.Split(address, ",")
	for _, addr := range addrs {
		hostPort := strings.Split(addr, ":")
		if len(hostPort) != 2 {
			return fmt.Errorf("Invalid address: %s", addr)
		}
		_, err := strconv.Atoi(hostPort[1])
		if err != nil {
			return err
		}
		gp.hosts = append(gp.hosts, addr)
	}
	gp.clients = make([]*nebula.Session, 0, concurrent)
	gp.channelBufferSize = chanSize
	gp.OutputCh = make(chan []string, gp.channelBufferSize)
	return nil
}

// Deprecated ConfigCsvStrategy sets csv reader strategy
func (gp *GraphPool) ConfigCsvStrategy(strategy int) {
	return
}

// ConfigCSV makes the read csv file configuration
func (gp *GraphPool) ConfigCSV(path, delimiter string, withHeader bool, opts ...interface{}) error {
	var (
		limit int = 500 * 10000
	)
	if gp.csvReader != nil {
		return nil
	}
	if len(opts) > 0 {
		l, ok := opts[0].(int)
		if ok {
			limit = l
		}
	}
	gp.csvReader = common.NewCsvReader(path, delimiter, withHeader, limit)

	if err := gp.csvReader.ReadForever(gp.DataCh); err != nil {
		return err
	}

	return nil
}

// ConfigOutput makes the output file configuration, would write the execution outputs
func (gp *GraphPool) ConfigOutput(path string) error {
	writer := common.NewCsvWriter(path, ",", outputHeader, gp.OutputCh)
	if err := writer.WriteForever(); err != nil {
		return err
	}
	return nil
}

// Close closes the nebula pool
func (gp *GraphPool) Close() error {
	gp.mutex.Lock()
	defer gp.mutex.Unlock()
	if !gp.initialized {
		return nil
	}
	// gp.Log.Println("begin close the nebula pool")
	for _, s := range gp.clients {
		if s != nil {
			s.Release()
		}
	}
	gp.initialized = false
	return nil
}

// GetSession gets the session from pool
func (gp *GraphPool) GetSession(username, password string) (common.IGraphClient, error) {
	gp.mutex.Lock()
	defer gp.mutex.Unlock()
	// balancer, ccore just use the first endpoint
	index := len(gp.clients) % len(gp.hosts)
	client, err := gp.clientGetter(gp.hosts[index], username, password)

	if err != nil {
		return nil, err
	}

	gp.clients = append(gp.clients, client)
	s := &GraphClient{Client: client, Pool: gp, DataCh: gp.DataCh}

	return s, nil
}

func (gc *GraphClient) Open() error {
	return nil
}
func (gc *GraphClient) Close() error {
	gc.Client.Release()
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
		resp       *nebula.ResultSet
		rows       int32
		latency    int64
	)
	start := time.Now()
	resp, err = gc.Client.Execute(stmt)

	if err != nil {
		isSucceed = false
		errMessage = err.Error()
		rows = 0
		latency = 0
	} else if !resp.IsSucceed() {
		isSucceed = false
		errMessage = resp.GetStatus()
		rows = 0
		latency = 0
	} else {
		rows = int32(resp.GetRowSize())
		latency = resp.GetLatency()
	}

	responseTime := int32(time.Since(start) / 1000)
	// output
	if gc.Pool.OutputCh != nil {
		var fr []string
		cols := resp.GetColSize()
		// print the first row of the result
		if rows != 0 {
			r, err := resp.GetRowValuesByIndex(0)
			if err != nil {
				return nil, err
			}
			for i := 0; i < cols; i++ {
				s, err := r.GetValueByIndex(i)
				if err != nil {
					return nil, err
				}
				fr = append(fr, s.String())
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
	if r.err != nil || !r.ResultSet.IsSucceed() {
		return false
	}

	return true
}

// GetLatency GetLatency
func (r *Response) GetLatency() int64 {
	if r.ResultSet != nil {
		return r.ResultSet.GetLatency()
	}
	return 0
}

// GetRowSize GetRowSize
func (r *Response) GetRowSize() int32 {
	if r.ResultSet != nil {
		return int32(r.ResultSet.GetRowSize())
	}
	return 0
}
