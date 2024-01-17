package cache

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type MetaSession struct {
	Address string
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

func SaveMetaSession(addr string) error {
	data, err := json.Marshal(MetaSession{Address: addr})
	if err != nil {
		return err
	}
	cacheFile := cachePath()
	err = ioutil.WriteFile(cacheFile, data, 0644)
	if err != nil {
		fmt.Println("Save meta session failed: ", err.Error())
		return err
	}
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
