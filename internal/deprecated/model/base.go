package model

// import (
// 	"context"
// 	"errors"
// 	"github.com/google/uuid"
// 	"log"
// 	"mals/pkg/config"
// )

// type ModelRequest struct {
// 	Text              string
// 	Schema            any
// 	SchemaName        string
// 	SchemaDescription string
// }

// type ModelTask struct {
// 	id      uuid.UUID
// 	Request ModelRequest
// 	Output  chan<- ModelResponse
// }

// type ModelResponse struct {
// 	Text  string
// 	Error error
// }

// type ModelService interface {
// 	Serve(ctx context.Context)
// 	onRequest(ctx context.Context, request ModelRequest) ModelResponse
// 	SubmitTask(request ModelTask)
// }

// type Model struct {
// 	ModelService
// 	logger   *log.Logger
// 	tasks    chan ModelTask
// 	Id       string
// 	Spec     string
// 	BaseUrl  string
// 	Settings config.ModelSettings
// }

// func NewModelResponse(text string) ModelResponse {
// 	return ModelResponse{Text: text, Error: nil}
// }

// func NewModelError(error error) ModelResponse {
// 	return ModelResponse{Text: "", Error: error}
// }

// func NewModelRequest(text string, schema any, schemaName string, schemaDescription string) ModelRequest {
// 	return ModelRequest{
// 		Text:              text,
// 		Schema:            schema,
// 		SchemaName:        schemaName,
// 		SchemaDescription: schemaDescription,
// 	}
// }

// // just in case there would be different kinds of requests
// func NewModelTask(request ModelRequest, output chan<- ModelResponse) ModelTask {
// 	return ModelTask{
// 		id:      uuid.New(),
// 		Request: request,
// 		Output:  output,
// 	}
// }

// func NewModel(logger *log.Logger, id string, spec string, baseUrl string, settings config.ModelSettings) Model {
// 	return Model{
// 		logger:   logger,
// 		tasks:    make(chan ModelTask),
// 		Id:       id,
// 		Spec:     spec,
// 		BaseUrl:  baseUrl,
// 		Settings: settings,
// 	}
// }

// func (m *Model) Serve(ctx context.Context) {
// 	m.LogInfoPrintf("serving started")
// 	for {
// 		select {
// 		case <-ctx.Done():
// 			close(m.tasks)
// 			m.LogInfoPrintf("serving stopped")
// 			return
// 		case task := <-m.tasks:
// 			m.LogInfoPrintf("task %s received", task.id)
// 			response := m.ModelService.onRequest(ctx, task.Request)
// 			m.LogInfoPrintf("task %s send", task.id)
// 			task.Output <- response
// 			m.LogInfoPrintf("task %s acknowledged", task.id)
// 		}
// 	}
// }

// func (m *Model) onRequest(ctx context.Context, request ModelRequest) ModelResponse {
// 	return NewModelError(errors.New("generic model is unable to generate responses"))
// }

// func (m *Model) SubmitTask(task ModelTask) {
// 	m.tasks <- task
// }
