package config

import (
	"encoding/json"
	"os"
	"sync"
)

const defaultComponentConfigPath = "etc/component-config.json"

var configMutex sync.RWMutex

type Component struct {
	Config map[string]interface{} `json:"config"`
}

type ComponentConfig struct {
	Graphd   Component `json:"graphd"`
	Metad    Component `json:"metad"`
	Storaged Component `json:"storaged"`
}

func LoadConfig() (*ComponentConfig, error) {
	configMutex.RLock()
	defer configMutex.RUnlock()

	data, err := os.ReadFile(defaultComponentConfigPath)
	if err != nil {
		if os.IsNotExist(err) {
			// If the file does not exist, return a new instance
			return &ComponentConfig{
				Graphd:   Component{Config: make(map[string]interface{})},
				Metad:    Component{Config: make(map[string]interface{})},
				Storaged: Component{Config: make(map[string]interface{})},
			}, nil
		}
		return nil, err
	}

	var config ComponentConfig
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func SaveConfig(config *ComponentConfig) error {
	configMutex.Lock()
	defer configMutex.Unlock()

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	// If the file does not exist, create a new one
	err = os.WriteFile(defaultComponentConfigPath, data, 0o644)
	if err != nil {
		return err
	}

	return nil
}

func (c *ComponentConfig) GetComponent(componentName string) *Component {
	switch componentName {
	case "graphd":
		return &c.Graphd
	case "metad":
		return &c.Metad
	case "storaged":
		return &c.Storaged
	default:
		return nil
	}
}

func (c *ComponentConfig) SetConfig(componentName string, config map[string]interface{}) {
	component := c.GetComponent(componentName)
	if component == nil {
		return
	}

	for k, v := range config {
		component.Config[k] = v
	}
}
