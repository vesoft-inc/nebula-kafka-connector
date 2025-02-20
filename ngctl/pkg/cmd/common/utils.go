package common

func GetResourceName(args []string) (string, error) {
	if len(args) < 1 {
		return "", NgctlError("Must specify a resource name", "")
	}
	return args[0], nil
}
