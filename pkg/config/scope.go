package config

type ScopeKind string

const (
	ScopeGlobal   ScopeKind = "global"
	ScopeClient   ScopeKind = "client"
	ScopeWorkflow ScopeKind = "workflow"
	ScopeStep     ScopeKind = "step"
)
