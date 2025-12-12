package model

func (s *ModelController) handleShutdown(t *TaskShutdown) {
	defer close(t.Result)

}

func (s *ModelController) handleRegister(t *TaskRegister) {
	defer close(t.Result)

}

func (s *ModelController) handleUnregister(t *TaskUnregister) {
	defer close(t.Result)

}

func (s *ModelController) handleCreate(t *TaskCreate) {
	defer close(t.Result)

}

func (s *ModelController) handleDelete(t *TaskDelete) {
	defer close(t.Result)

}

func (s *ModelController) handleStart(t *TaskStart) {
	defer close(t.Result)

}

func (s *ModelController) handleStop(t *TaskStop) {
	defer close(t.Result)

}
