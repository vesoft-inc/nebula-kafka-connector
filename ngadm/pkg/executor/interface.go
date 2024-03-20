package executor

type Executor interface {
	Shell(cmd string, sudo bool) (stdout string, stderr string, err error)
	Upload(src string, dest string) (stdout string, stderr string, err error)
	Download(src string, dest string) (stdout string, stderr string, err error)
}
