package flistfile

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	flistfileName = "flist.txt"
)

// FListFile 構造体
type FListFile struct {
	// Dir 走査対象ディレクトリ
	Dir string
	// flistファイルの出力先パス
	OutPath string
}

// FileRowProp tmpファイル1行の情報を保持する構造体
type FileRowProp struct {
	Path string
	Name string
}

// NewFListFile はFListFileのコンストラクタ
func NewFListFile(dir string) (*FListFile, error) {
	ep, err := os.Executable()
	if err != nil {
		return nil, err
	}
	fp := filepath.Join(filepath.Dir(ep), flistfileName)

	return &FListFile{
		Dir:     dir,
		OutPath: fp,
	}, nil
}

func getNewLineCode() string {
	if runtime.GOOS == "windows" {
		return "\r\n"
	}
	return "\n"
}

// Create flistファイルを作成する。
// 返り値に走査対象のディレクトリにあるファイル一覧情報を返却する
func (f *FListFile) Create() ([]FileRowProp, error) {
	file, err := os.OpenFile(f.OutPath, os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// 引数Directoryからファイル一覧読み込み
	rows, err := readFileRowProps(f.Dir)
	if err != nil {
		return nil, err
	}

	nlc := getNewLineCode()

	// ファイル一覧の書き込み
	for _, p := range rows {
		_, err := file.WriteString(p.Name + nlc)
		if err != nil {
			return nil, err
		}
	}
	return rows, nil
}

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

// OpenWithEditor flistファイルを指定エディタで開く
func (f *FListFile) OpenWithEditor(name string) error {
	execCmd := exec.Command(name, f.OutPath)
	return execCmd.Run()
}

// OpenRead はflistファイルを読み込みProperty情報を返却する
func (f *FListFile) OpenRead() ([]FileRowProp, error) {
	var props []FileRowProp

	file, err := os.Open(f.OutPath)
	if err != nil {
		return props, err
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return props, err
	}

	rows := strings.Split(string(b), getNewLineCode())

	for _, r := range rows {
		props = append(props, FileRowProp{
			Path: filepath.Join(f.Dir, r),
			Name: r,
		})
	}

	return props, nil
}

// Remove flistファイルを削除する
func (f *FListFile) Remove() error {
	return os.Remove(f.OutPath)
}
