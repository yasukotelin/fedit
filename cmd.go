package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yasukotelin/fedit/editor"
	"github.com/yasukotelin/fedit/file"
)

var (
	specifiedEditor string

	rootCmd = &cobra.Command{
		Use:   "fedit",
		Short: "rename all files in derectory.",
		Long:  "fedit is the tool to rename all files in directory",
		Run:   run,
	}
)

func init() {
	rootCmd.Version = version
	rootCmd.PersistentFlags().StringVarP(&specifiedEditor, "editor", "e", editor.GetDefaultEditor(), "specify the editor to open. ")
}

func exitError(e error) {
	fmt.Println(e)
	os.Exit(1)
}

func exitErrorS(s string) {
	fmt.Println(s)
	os.Exit(1)
}

func run(cmd *cobra.Command, args []string) {
	dirPath := getDirPath(args)

	flFile, err := file.NewTmpFile(tmpFileName, dirPath)
	if err != nil {
		exitError(err)
	}

	// tmpファイル作成
	f1, err := flFile.Create()
	if err != nil {
		exitError(err)
	}

	// エディタでtmpファイルを開く
	if err := flFile.OpenWithEditor(specifiedEditor); err != nil {
		exitError(err)
	}

	// 編集後のtmpファイルを開いて読み込む
	f2, err := flFile.OpenRead()
	if err != nil {
		exitError(err)
	}

	// リネーム名に重複がないかのチェック
	if file.IsDupl(&f2) {
		exitError(errors.New("Duplicate file path specified"))
	}

	// 差分取得
	diffs := file.Diff(f1, f2)

	if len(diffs) > 0 {
		// 差分表示
		fmt.Println()
		for _, d := range diffs {
			fmt.Printf("%s ---> %s\n", d.File1.Path, d.File2.Path)
		}
		fmt.Println()

		// 確定確認
		r, err := askToApplyRename()
		if err != nil {
			exitError(err)
		}
		if r {
			// Rename処理
			for _, d := range diffs {
				if err := os.Rename(d.File1.Path, d.File2.Path); err != nil {
					exitError(err)
				}
			}
		}
	}

	// tmpファイル削除
	if err := flFile.Remove(); err != nil {
		exitError(err)
	}
}

func getDirPath(args []string) string {
	if len(args) == 0 {
		exitErrorS("directry path required")
	}
	return args[0]
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
