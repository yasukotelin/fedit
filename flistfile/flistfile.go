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
	FlistfileName = "flist.txt"
)

// Flistfile 構造体
type Flistfile struct {
	// Dir 走査対象ディレクトリ
	Dir string
	// flistファイルの出力先パス
	OutPath string
}

// NewFlistfile はFlistfileのコンストラクタ
func NewFlistfile(dir string) (*Flistfile, error) {
	ep, err := os.Executable()
	if err != nil {
		return nil, err
	}
	fp := filepath.Join(filepath.Dir(ep), FlistfileName)

	return &Flistfile{
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
func (f *Flistfile) Create() ([]Row, error) {
	file, err := os.OpenFile(f.OutPath, os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// 引数Directoryからファイル一覧読み込み
	rows, err := readRows(f.Dir)
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
func readRows(dir string) ([]Row, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var res []Row
	index := 0
	for _, f := range files {
		if !f.IsDir() {
			row := Row{
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
func (f *Flistfile) OpenWithEditor(name string) error {
	execCmd := exec.Command(name, f.OutPath)
	return execCmd.Run()
}

// OpenRead はflistファイルを読み込みProperty情報を返却する
func (f *Flistfile) OpenRead() ([]Row, error) {
	var props []Row

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
		props = append(props, Row{
			Path: filepath.Join(f.Dir, r),
			Name: r,
		})
	}

	return props, nil
}

// Remove flistファイルを削除する
func (f *Flistfile) Remove() error {
	return os.Remove(f.OutPath)
}
