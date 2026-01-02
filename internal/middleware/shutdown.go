package middleware

func (s *Middleware) Shutdown() error {
	s.workspaces.Range(func(key string, value *Workspace) bool {
		s.workspaceDelete(key)
		return true
	})

	s.initialized = false

	return nil
}
