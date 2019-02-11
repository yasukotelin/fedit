package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yasukotelin/feditelin/tmpfile"
)

const (
	version = "0.1.0"
)

var (
	argEditor string

	rootCmd = &cobra.Command{
		Use:   "feditelin",
		Short: "rename all files in derectory.",
		Long:  "feditelin is the tool to rename all files in directory",
		Run:   run,
	}
)

func init() {
	rootCmd.Version = version

	// Flag定義
	rootCmd.PersistentFlags().StringVarP(&argEditor, "editor", "e", "", "specify the editor to open. ")
}

func run(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println("directry path required")
		os.Exit(1)
	}

	tmpFile := tmpfile.NewTmpFile(args[0])
	if err := tmpFile.Create(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := tmpFile.OpenWithEditor(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := tmpFile.Delete(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
