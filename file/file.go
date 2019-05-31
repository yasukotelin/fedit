package file

import "runtime"

func getNewLineCode() string {
	if runtime.GOOS == "windows" {
		return "\r\n"
	}
	return "\n"
}
