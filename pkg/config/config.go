package config

type Config struct {
	Logs      []*Log
	Models    []*Model
	Lsps      []*Lsp
	Handlers  []*Handler
	Listeners []*Listener
}
