package cmd

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
	f1, err := flistfile.Create(dirPath, fPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// エディタでtmpファイルを開く
	execCmd := exec.Command(editor, fPath)
	if err := execCmd.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// 編集後のtmpファイルを開いて読み込む
	f2, err := flistfile.OpenRead(fPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// 差分取得
	diffs := flistfile.Diff(f1, f2)

	if len(diffs) > 0 {
		// 差分表示
		fmt.Println()
		for _, d := range diffs {
			fmt.Printf("%s ---> %s\n", d.File1.Name, d.File2.Name)
		}
		fmt.Println()

		// 確定確認
		r, err := askToApplyRename()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if r {
			// Rename処理
			for _, d := range diffs {
				if err := os.Rename(d.File1.Path, d.File2.Path); err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
			}
		}
	}

	// tmpファイル削除
	if err := os.Remove(fPath); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func askToApplyRename() (bool, error) {
	var s string
	for {
		fmt.Print("apply to rename [y/n]? ")
		_, err := fmt.Scanln(&s)
		if err != nil {
			return false, err
		}
		switch s {
		case "y", "Y":
			return true, nil
		case "n", "N":
			return false, nil
		}
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
