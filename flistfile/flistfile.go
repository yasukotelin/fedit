package flistfile

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// const (
// 	fname = "fed-tmp.txt"
// )

var (
	descriMsg = []string{
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
)

// // TmpFile struct
// type TmpFile struct {
// 	Dir         string
// 	Name        string
// 	Rows        []FileRowProp
// 	NewLineCode string
// 	Editor      string
// }

// FileRowProp tmpファイル1行の情報を保持する構造体
type FileRowProp struct {
	Path string
	Name string
}

// FileRowDiff ファイル行の差分用構造体
type FileRowDiff struct {
	File1 FileRowProp
	File2 FileRowProp
}

func getNewLineCode() string {
	if runtime.GOOS == "windows" {
		return "\r\n"
	}
	return "\n"
}

// Create 指定dirパスにあるファイル一覧を書き込んだファイルをpathに作成する
// 返り値に読み込んだファイル一覧情報を返却する
func Create(dir string, path string) ([]FileRowProp, error) {
	file, err := os.OpenFile(path, os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// 引数Directoryからファイル一覧読み込み
	rows, err := readFileRowProps(dir)
	if err != nil {
		return nil, err
	}

	nlc := getNewLineCode()

	// ファイル説明文の書き込み
	file.WriteString(strings.Join(descriMsg, nlc) + nlc)

	// ファイル一覧の書き込み
	for _, p := range rows {
		_, err := file.WriteString(p.Name + nlc)
		if err != nil {
			return nil, err
		}
	}
	return rows, nil
}

// NewTmpFile create TmpFile instance
// func NewTmpFile(dir string) *TmpFile {
// 	var newLineCode string
// 	var editor string
// 	if runtime.GOOS == "windows" {
// 		editor = "notepad"
// 		newLineCode = "\r\n"
// 	} else {
// 		editor = "vim"
// 		newLineCode = "\n"
// 	}

// 	return &TmpFile{
// 		Dir:         dir,
// 		Name:        name,
// 		Rows:        nil,
// 		NewLineCode: newLineCode,
// 		Editor:      editor,
// 	}
// }

// func (t *TmpFile) getPath() string {
// 	return filepath.Join(t.Dir, t.Name)
// }

// func (t *TmpFile) Create() error {
// 	file, err := os.Create(t.getPath())
// 	if err != nil {
// 		return err
// 	}
// 	defer file.Close()

// 	// 引数Directoryからファイル一覧読み込み
// 	t.Rows, err = readFiles(t.Dir)
// 	if err != nil {
// 		return err
// 	}

// 	// ファイル説明文の書き込み
// 	file.WriteString(strings.Join(descriMsg, t.NewLineCode) + t.NewLineCode)

// 	// ファイル一覧の書き込み
// 	for _, p := range t.Rows {
// 		_, err := file.WriteString(p.Name + t.NewLineCode)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

// readFiles ディレクトリからファイル一覧を返却する
func readFileRowProps(dir string) ([]FileRowProp, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var res []FileRowProp
	index := 0
	for _, f := range files {
		if !f.IsDir() {
			row := FileRowProp{
				Path: filepath.Join(dir, f.Name()),
				Name: f.Name(),
			}
			res = append(res, row)
			index++
		}
	}
	return res, nil
}

// func (t *TmpFile) OpenWithEditor() error {
// 	execCmd := exec.Command(t.Editor, t.getPath())
// 	if err := execCmd.Run(); err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (t *TmpFile) Delete() error {
// 	return os.Remove(t.getPath())
// }

// TODO implements
// func Diff(f1 *TmpFile, f2 *TmpFile) ([]FileRowDiff, error) {
// 	var diffs []FileRowDiff
// 	return diffs, nil
// }

// TODO implements
// func (t *TmpFile) OpenRead() (string, error) {
// 	f, err := os.Open(t.getPath())
// 	if err != nil {
// 		return "", err
// 	}

// 	b, err := ioutil.ReadAll(f)
// 	if err != nil {
// 		return "", err
// 	}

// 	return string(b), nil
// }
