package types

type Logger interface {
	Info(msg string)
	Warn(msg string)
	Error(msg string)
	Fatal(msg string)
}

type Progress struct {
	Total        int      `json:"total,omitempty"`
	Current      int      `json:"current,omitempty"`
	Ratio        float32  `json:"ratio,omitempty"`
	RunningTasks []string `json:"runningTasks,omitempty"` // maybe we can show each task's progress
}
