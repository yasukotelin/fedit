package flistfile

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

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

func OpenRead(path string) ([]FileRowProp, error) {
	var props []FileRowProp
	dir := filepath.Dir(path)

	f, err := os.Open(path)
	if err != nil {
		return props, err
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return props, err
	}

	rows := strings.Split(string(b), getNewLineCode())

	for _, r := range rows {
		props = append(props, FileRowProp{
			Path: filepath.Join(dir, r),
			Name: r,
		})
	}

	return props, nil
}
