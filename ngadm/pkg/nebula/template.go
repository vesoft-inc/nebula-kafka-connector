package nebula

import (
	"bytes"

	"github.com/unknwon/goconfig"
)

type Template struct {
	*goconfig.ConfigFile
}

func NewTemplateWithData(data []byte) (*Template, error) {
	cfg, err := goconfig.LoadFromData(data)
	if err != nil {
		return nil, err
	}
	cfg.SetPrettyFormat(false)
	return &Template{cfg}, err
}

func (t *Template) AddReferConfig(referConfig *Template) {
	referKeys := referConfig.ConfigFile.GetKeyList("")
	curKeyMap, _ := t.ConfigFile.GetSection("")
	for _, key := range referKeys {
		if _, ok := curKeyMap[key]; ok {
			value, _ := referConfig.ConfigFile.GetValue("", key)
			t.ConfigFile.SetValue("", key, value)
		}
	}
}

func (t *Template) UpdateConfig(config map[string]string) {
	for k, v := range config {
		t.ConfigFile.SetValue("", k, v)
	}
}

func (t *Template) DeleteKeys(keys []string) {
	for _, key := range keys {
		t.ConfigFile.DeleteKey("", key)
	}
}

func (t *Template) String() (string, error) {
	var data bytes.Buffer
	if err := goconfig.SaveConfigData(t.ConfigFile, &data); err != nil {
		return "", err
	}
	return data.String(), nil
}
