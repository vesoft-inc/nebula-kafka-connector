package types

type WorkflowSpec struct {
	Rollback    bool        `json:"rollback,omitempty"`
	Type        string      `json:"type,omitempty"` // default tasks is indeed a serial task
	Params      any         `json:"params,omitempty"`
	Description string      `json:"description,omitempty"`
	Tasks       []*TaskSpec `json:"tasks,omitempty"`
}

type TaskSpec struct {
	Type        string      `json:"type,omitempty"`
	Description string      `json:"description,omitempty"`
	Params      any         `json:"params,omitempty"`
	SubTasks    []*TaskSpec `json:"subTasks,omitempty"`
}
