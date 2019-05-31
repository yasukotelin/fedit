package file

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// TmpFile is temporary file used for editing it and confirming diff.
type TmpFile struct {
	// Name is tmporary file name
	Name string
	// Dir is target directory to search files.
	Dir string
	// OutPath is path where this is output
	OutPath string
}

// NewTmpFile returns created struct pointer
func NewTmpFile(name string, dir string) (*TmpFile, error) {
	exePath, err := os.Executable()
	if err != nil {
		return nil, err
	}
	outPath := filepath.Join(filepath.Dir(exePath), name)

	return &TmpFile{
		Name:    name,
		Dir:     dir,
		OutPath: outPath,
	}, nil
}

// Create flistファイルを作成する。
// 返り値に走査対象のディレクトリにあるファイル一覧情報を返却する
func (f *TmpFile) Create() ([]Row, error) {
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
func (f *TmpFile) OpenWithEditor(name string) error {
	execCmd := exec.Command(name, f.OutPath)
	return execCmd.Run()
}

// OpenRead はflistファイルを読み込みProperty情報を返却する
func (f *TmpFile) OpenRead() ([]Row, error) {
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
func (f *TmpFile) Remove() error {
	return os.Remove(f.OutPath)
}
