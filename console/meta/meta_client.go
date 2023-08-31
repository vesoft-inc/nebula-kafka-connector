package meta

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	nrpc "github.com/vesoft-inc/nebula-ng-tools/golang/nrpc"
)

type MetaClient struct {
	client *nrpc.Client
}

type MetaSession struct {
	Address       string
	LeaderAddress string
}

type LeaderAddress struct {
	Host string
	Port uint32
}

func cachePath() string {
	cacheHome := os.Getenv("HOME")
	if cacheHome == "" {
		ex, err := os.Executable()
		if err != nil {
			log.Panicf("Get executable failed: %s", err.Error())
		}
		cacheHome = filepath.Dir(ex) // Set to executable folder
	}
	cacheFile := filepath.Join(cacheHome, ".nebula_meta_session")
	return cacheFile
}

func SaveMetaSession(addr string, leader string) error {
	data, err := json.Marshal(MetaSession{Address: addr, LeaderAddress: leader})
	if err != nil {
		return err
	}
	cacheFile := cachePath()
	err = ioutil.WriteFile(cacheFile, data, 0644)
	if err != nil {
		fmt.Println("Save meta session failed: ", err.Error())
		return err
	}
	fmt.Println("Saving leader meta:", leader)
	return nil
}

func LoadMetaSession() (*MetaSession, error) {
	cachePath := cachePath()
	data, err := ioutil.ReadFile(cachePath)
	if err != nil {
		return nil, err
	}
	var metaSession MetaSession
	err = json.Unmarshal(data, &metaSession)
	if err != nil {
		return nil, err
	}
	return &metaSession, nil
}

func NewMetaClient(addr string) *MetaClient {
	metaclient := &MetaClient{
		client: nrpc.NewClient(addr),
	}
	return metaclient
}

func (m *MetaClient) Close() {
	m.client.Close()
}

func LoadMetaClient() (*MetaClient, error) {
	session, err := LoadMetaSession()
	if err != nil {
		return nil, err
	}
	meta := NewMetaClient(session.LeaderAddress)
	leader, err := meta.GetLeader(NewGetLeaderRequest())
	if err != nil {
		metas := strings.Split(session.Address, ",")
		if len(metas) == 0 {
			return nil, fmt.Errorf("meta server address is empty")
		}
		for _, m := range metas {
			fmt.Println("Try to connect to meta:", m)
			metaclient := NewMetaClient(m)
			leader, err = metaclient.GetLeader(NewGetLeaderRequest())
			if err != nil {
				fmt.Println("get leader failed: ", err)
				continue
			}
			break
		}
		if err != nil {
			return nil, fmt.Errorf("All get leader failed: %s", err.Error())
		}
	}
	leaderStr := fmt.Sprintf("%s:%d", leader.Host, leader.Port)

	if leaderStr != session.LeaderAddress {
		fmt.Println("Leader changed, reconnecting...")
		meta = NewMetaClient(fmt.Sprintf("%s:%d", leader.Host, leader.Port))
		SaveMetaSession(session.Address, leaderStr)
	}
	fmt.Println("Using leader meta:", leaderStr)

	return meta, nil
}

func (c *MetaClient) Login(user string, password string) error {
	return nil
}

func (c *MetaClient) GetLeader(request *GetLeaderRequest) (LeaderAddress, error) {
	bytes := request.Serialize()
	resp, err := c.send(bytes)
	if err != nil {
		return LeaderAddress{"", 0}, err
	}
	deserializer := NewDeserializer(resp)
	respHeader := DeserializeHeader(deserializer)
	if respHeader.code != 0 {
		return LeaderAddress{"", 0}, fmt.Errorf("Error: code: %d msg:%s", respHeader.code, respHeader.msg)
	}
	getLeaderResp := DeserializeGetLeaderResponse(deserializer)

	return LeaderAddress{getLeaderResp.host, getLeaderResp.port}, nil
}

func (c *MetaClient) CreateCluster(request *CreateClusterRequest) (string, time.Duration, error) {
	bytes := request.Serialize()
	start := time.Now()
	resp, err := c.send(bytes)
	duration := time.Since(start)
	if err != nil {
		return err.Error(), duration, err
	}

	deserializer := NewDeserializer(resp)
	respHeader := DeserializeHeader(deserializer)
	if respHeader.code != 0 {
		return fmt.Sprintf("Error: code: %d msg:%s", respHeader.code, respHeader.msg), duration, nil
	}

	return "success", duration, nil
}

func (c *MetaClient) AddService(request *AddServiceRequest) (string, time.Duration, error) {
	bytes := request.Serialize()
	start := time.Now()
	resp, err := c.send(bytes)
	duration := time.Since(start)
	if err != nil {
		return err.Error(), duration, err
	}
	deserializer := NewDeserializer(resp)
	respHeader := DeserializeHeader(deserializer)
	if respHeader.code != 0 {
		return fmt.Sprintf("Error: code: %d msg:%s", respHeader.code, respHeader.msg), duration, nil
	}

	return "success", duration, nil
}

func (c *MetaClient) InitCluster(request *InitClusterRequest) (string, time.Duration, error) {
	bytes := request.Serialize()
	start := time.Now()
	resp, err := c.send(bytes)
	duration := time.Since(start)
	if err != nil {
		return "", 0, err
	}
	deserializer := NewDeserializer(resp)
	respHeader := DeserializeHeader(deserializer)
	if respHeader.code != 0 {
		return fmt.Sprintf("Error: code: %d msg:%s", respHeader.code, respHeader.msg), duration, nil
	}

	return "success", duration, nil
}

func (c *MetaClient) CreateSchema(request *MetaDDLRequest) (string, time.Duration, error) {
	bytes := request.Serialize()
	start := time.Now()
	resp, err := c.send(bytes)
	duration := time.Since(start)
	if err != nil {
		return "", 0, err
	}
	deserializer := NewDeserializer(resp)
	respHeader := DeserializeHeader(deserializer)
	if respHeader.code != 0 {
		return fmt.Sprintf("Error: code: %d msg:%s", respHeader.code, respHeader.msg), duration, nil
	}

	return "success", duration, nil
}

func (c *MetaClient) ShowService(request *ListServiceRequest) (string, time.Duration, error) {
	bytes := request.Serialize()
	start := time.Now()
	resp, err := c.send(bytes)
	duration := time.Since(start)
	if err != nil {
		return "", 0, err
	}
	deserializer := NewDeserializer(resp)
	respHeader := DeserializeHeader(deserializer)
	if respHeader.code != 0 {
		return fmt.Sprintf("Error: code: %d msg:%s", respHeader.code, respHeader.msg), duration, nil
	}

	respBody := DeserializeListServiceResponse(deserializer)

	return respBody.Format(), duration, nil
}

func (c *MetaClient) ShowCluster(request *ListClusterRequest) (string, time.Duration, error) {
	bytes := request.Serialize()
	start := time.Now()
	resp, err := c.send(bytes)
	duration := time.Since(start)
	if err != nil {
		return "", 0, err
	}
	deserializer := NewDeserializer(resp)
	respHeader := DeserializeHeader(deserializer)
	if respHeader.code != 0 {
		return fmt.Sprintf("Error: code: %d msg:%s", respHeader.code, respHeader.msg), duration, nil
	}

	respBody := DeserializeListClusterResponse(deserializer)

	return respBody.Format(), duration, nil
}

func (c *MetaClient) send(request []byte) ([]byte, error) {
	timeout := 200 * time.Millisecond
	resp, err := c.client.Send(request, timeout)
	if err != nil {
		nerr, ok := err.(nrpc.Error)
		if ok {
			if nerr.Timeout() {
				fmt.Println("Timeout")
				return nil, err
				// 'client' is still healthy to use
			} else if nerr.BadChannel() {
				fmt.Printf("BadChannel: ")
				fmt.Println("Reconnecting...")
				if err = c.client.Reconnect(3, time.Second); err != nil {
					fmt.Println(err)
					return nil, err
				}
				resp, err := c.client.Send(request, timeout)
				if err != nil {
					fmt.Println(err)
					return nil, err
				}
				return resp, nil
			} else {
				fmt.Println("Unknown rpc error")
				return nil, err
			}
		}
	}
	return resp, nil
}
