package main

import (
	"fmt"
	"os"
)

const (
	version     = "1.1.3"
	tmpFileName = "files.tmp"
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(ExitCodeError)
	}
	os.Exit(ExitCodeOk)
}
