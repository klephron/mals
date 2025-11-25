package common

type LogGroup struct {
	logs []Log
}

func NewGroup(logs []Log) *LogGroup {
	return &LogGroup{
		logs: logs,
	}
}

func (g *LogGroup) Debug(msg string, args ...any) {
	for _, log := range g.logs {
		log.Debug(msg, args...)
	}
}

func (g *LogGroup) Info(msg string, args ...any) {
	for _, log := range g.logs {
		log.Info(msg, args...)
	}
}

func (g *LogGroup) Warn(msg string, args ...any) {
	for _, log := range g.logs {
		log.Warn(msg, args...)
	}
}

func (g *LogGroup) Error(msg string, args ...any) {
	for _, log := range g.logs {
		log.Error(msg, args...)
	}
}
