package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

const (
	version     = "0.1.0"
	tmpFileName = "tmp-file-list.txt"
)

var (
	editor      string
	defEditor   string
	newLineCode string

	rootCmd = &cobra.Command{
		Use:   "feditelin",
		Short: "rename all files in derectory.",
		Long:  "feditelin is the tool to rename all files in directory",
		Run:   run,
	}
)

func init() {
	if runtime.GOOS == "windows" {
		defEditor = "notepad"
		newLineCode = "\r\n"
	} else {
		defEditor = "vim"
		newLineCode = "\n"
	}

	rootCmd.Version = version

	// Flag定義
	rootCmd.PersistentFlags().StringVarP(&editor, "editor", "e", defEditor, "specify the editor to open. ")
}

func run(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println("directry path required")
		os.Exit(1)
	}

	tmpFilePath, err := createTmpFile(args[0])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := openEditor(editor, tmpFilePath); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// 後始末（tmpファイルの削除）
	err = os.Remove(tmpFilePath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// Execute root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func createTmpFile(dir string) (tmpFilePath string, err error) {
	exeDir, err := getExeDirPath()
	if err != nil {
		return tmpFileName, err
	}

	tmpFilePath = filepath.Join(exeDir, tmpFileName)

	file, err := os.Create(tmpFilePath)
	if err != nil {
		return tmpFileName, err
	}
	defer file.Close()

	// 引数Directoryからファイル一覧読み込み
	fileProps, err := readFiles(dir)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	file.WriteString(makeTmpFileDescriMsg())

	for _, p := range fileProps {
		_, err := file.WriteString(p.Name + newLineCode)
		if err != nil {
			return tmpFileName, err
		}
	}

	return tmpFileName, nil
}

// OpenEditor 指定editorでfPathを実行する
func openEditor(editor string, fPath string) error {
	execCmd := exec.Command(editor, fPath)
	if err := execCmd.Run(); err != nil {
		return err
	}
	return nil
}

// GetExeDirPath アプリ実行ディレクトリのパスを返却する
func getExeDirPath() (string, error) {
	p, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Dir(p), nil
}

// ReadFiles ディレクトリからファイル一覧を返却する
func readFiles(dir string) ([]FileRowProp, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var res []FileRowProp
	index := 0
	for _, f := range files {
		if !f.IsDir() {
			row := FileRowProp{
				Index:   index,
				Path:    filepath.Join(dir, f.Name()),
				Name:    f.Name(),
				NewName: "",
			}
			res = append(res, row)
			index++
		}
	}
	return res, nil
}

func makeTmpFileDescriMsg() string {
	tmpFileDescriMsg := []string{
		"// リネームしたいファイル名を編集してください。",
		"// ファイル名の編集が完了したら保存して閉じてください。",
		"// 編集したファイル名へのリネーム処理が実行されます。",
		"",
		"// ファイルのリネーム判定は表示されている既定の行に対して行われます。",
		"// 改行や行削除などを実施し行をずらした場合には、",
		"// ファイル名が一様に変更されてしまうため注意してください。",
		"",
		"// 編集が完了し閉じた後、",
		"// 最終確認が表示されるのでキャンセルしたい場合は一度保存して次へ進めてください。",
		"",
		"// ----- files ------",
	}
	return strings.Join(tmpFileDescriMsg, newLineCode) + newLineCode
}
