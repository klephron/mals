package config

import "encoding/json"

type Model struct {
	Id       string `json:"id"`
	BaseUrl  string `json:"base_url"`
	Spec     string `json:"spec"`
	Settings any    `json:"settings"`
}

type LspServer struct {
	Name      string   `json:"name"`
	Filetypes []string `json:"filetypes"`
	Cmd       []string `json:"cmd"`
	Settings  any      `json:"settings"`
}

type Workspace struct {
	LspServers []LspServer `json:"lsp_servers"`
	Model      Model       `json:"model"`
}

type Workspaces struct {
	Default Workspace `json:"*"`
}

type Config struct {
	Models     []Model    `json:"models"`
	Workspaces Workspaces `json:"workspaces"`
}

func (c *Config) String() string {
	bytes, _ := Encode(c)
	return string(bytes)
}

func Default() *Config {
	var config Config
	config.Workspaces.Default.LspServers = make([]LspServer, 0)
	config.Workspaces.Default.Model.Settings = struct{}{}

	config.Models = make([]Model, 0)
	return &config
}

func Decode(data []byte) (*Config, error) {
	config := Default()

	if err := json.Unmarshal(data, config); err != nil {
		return nil, err
	}

	for i, model := range config.Models {
		if model.Settings == nil {
			config.Models[i].Settings = struct{}{}
		}
	}

	for i, server := range config.Workspaces.Default.LspServers {
		if server.Settings == nil {
			config.Workspaces.Default.LspServers[i].Settings = struct{}{}
		}
	}

	return config, nil
}

func Encode(config *Config) ([]byte, error) {
	content, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}
	return content, nil
}
