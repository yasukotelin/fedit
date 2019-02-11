package tmpfile

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	name = "fed-tmp.txt"
)

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

// TmpFile struct
type TmpFile struct {
	Dir         string
	Name        string
	Rows        []FileRowProp
	NewLineCode string
	Editor      string
}

// FileRowProp tmpファイル1行の情報を保持する構造体
type FileRowProp struct {
	Index   int
	Path    string
	Name    string
	NewName string
}

// NewTmpFile create TmpFile instance
func NewTmpFile(dir string) *TmpFile {
	var newLineCode string
	var editor string
	if runtime.GOOS == "windows" {
		editor = "notepad"
		newLineCode = "\r\n"
	} else {
		editor = "vim"
		newLineCode = "\n"
	}

	return &TmpFile{
		Dir:         dir,
		Name:        name,
		Rows:        nil,
		NewLineCode: newLineCode,
		Editor:      editor,
	}
}

func (t *TmpFile) getPath() string {
	return filepath.Join(t.Dir, t.Name)
}

func (t *TmpFile) Create() error {
	file, err := os.Create(t.getPath())
	if err != nil {
		return err
	}
	defer file.Close()

	// 引数Directoryからファイル一覧読み込み
	t.Rows, err = readFiles(t.Dir)
	if err != nil {
		return err
	}

	// ファイル説明文の書き込み
	file.WriteString(strings.Join(descriMsg, t.NewLineCode) + t.NewLineCode)

	// ファイル一覧の書き込み
	for _, p := range t.Rows {
		_, err := file.WriteString(p.Name + t.NewLineCode)
		if err != nil {
			return err
		}
	}

	return nil
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

func (t *TmpFile) OpenWithEditor() error {
	execCmd := exec.Command(t.Editor, t.getPath())
	if err := execCmd.Run(); err != nil {
		return err
	}
	return nil
}

func (t *TmpFile) Delete() error {
	return os.Remove(t.getPath())
}
