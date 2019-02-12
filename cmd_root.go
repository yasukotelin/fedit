package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/yasukotelin/feditelin/flistfile"
)

const (
	version       = "0.1.0"
	flistfileName = "flist.txt"
)

var (
	editor    string
	defEditor string

	rootCmd = &cobra.Command{
		Use:   "feditelin",
		Short: "rename all files in derectory.",
		Long:  "feditelin is the tool to rename all files in directory",
		Run:   run,
	}
)

func init() {
	rootCmd.Version = version

	if runtime.GOOS == "windows" {
		defEditor = "notepad"
	} else {
		defEditor = "vim"
	}

	// Flag定義
	rootCmd.PersistentFlags().StringVarP(&editor, "editor", "e", defEditor, "specify the editor to open. ")
}

func run(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println("directry path required")
		os.Exit(1)
	}

	dirPath := args[0]
	fPath := filepath.Join(dirPath, flistfileName)

	// tmpファイル作成
	fProps, err := flistfile.Create(dirPath, fPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, prop := range fProps {
		fmt.Println(prop.Path)
	}

	// エディタでtmpファイルを開く
	execCmd := exec.Command(editor, fPath)
	if err := execCmd.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// tmpファイルを開いて読み込む
	flistfile.OpenRead(fPath)

	// tmpファイル削除
	if err := os.Remove(fPath); err != nil {
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
