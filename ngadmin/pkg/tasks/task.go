package tasks

type Task interface {
	Execute() error
	Rollback() error
	String() string
}
