package common

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

var CachePath string

type MetaToken struct {
	Address        string
	Leader         string
	Token          []byte
	EnableTLS      bool
	CA             string
	Cert           string
	Key            string
	PeerNameVerify bool
	PeerName       string
}

func GetCachePath() string {
	if CachePath != "" {
		return CachePath
	}
	cacheHome := os.Getenv("HOME")
	if cacheHome == "" {
		ex, err := os.Executable()
		if err != nil {
			log.Panicf("Get executable failed: %s", err.Error())
		}
		cacheHome = filepath.Dir(ex) // Set to executable folder
	}
	CachePath = filepath.Join(cacheHome, ".nebula_meta_token")

	return CachePath
}

func SaveMetaToken(addr string, leader string, token []byte, enableTLS bool, ca, cert, key string, peerNameVerify bool, peerName string) error {
	data, err := json.Marshal(MetaToken{Address: addr, Leader: leader, Token: token, EnableTLS: enableTLS, CA: ca, Cert: cert, Key: key, PeerNameVerify: peerNameVerify, PeerName: peerName})
	if err != nil {
		return err
	}
	cacheFile := GetCachePath()
	err = os.WriteFile(cacheFile, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func LoadMetaToken() (*MetaToken, error) {
	cachePath := GetCachePath()
	data, err := os.ReadFile(cachePath)
	if err != nil {
		return nil, err
	}
	var metaSession MetaToken
	err = json.Unmarshal(data, &metaSession)
	if err != nil {
		return nil, err
	}
	return &metaSession, nil
}

func ClearMetaToken() error {
	cachePath := GetCachePath()
	return os.Remove(cachePath)
}
