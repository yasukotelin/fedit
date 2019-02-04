package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

const (
	version     = "0.1.0"
	defEditor   = "vim"
	tmpFileName = "tmp-file-list.txt"
)

var (
	editor string

	rootCmd = &cobra.Command{
		Use:   "feditelin",
		Short: "rename all files in derectory.",
		Long:  "feditelin is the tool to rename all files in directory",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Println("directry path required")
				os.Exit(1)
			}

			fPath := filepath.Join(getExeDirPath(), tmpFileName)
			openEditor(editor, fPath)
		},
	}
)

func init() {
	rootCmd.Version = version

	// Flag定義
	rootCmd.PersistentFlags().StringVarP(&editor, "editor", "e", defEditor, "specify the editor to open. ")
}

// Execute root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func getExeDirPath() string {
	p, err := os.Executable()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return filepath.Dir(p)
}

func openEditor(editor string, fPath string) {
	execCmd := exec.Command(editor, fPath)
	if err := execCmd.Run(); err != nil {
		fmt.Println("Command Exec Error.")
		os.Exit(1)
	}
}
