package file

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

// TmpFile is temporary file used for editing it and confirming diff.
type TmpFile struct {
	// Name is tmporary file name
	Name string
	// TargetDir is target directory to search files.
	TargetDir string
	// OutPath is path where this is output
	OutPath string
	// Rows is file rows
	Rows []Row
	// PrevRows is row before edited (previous rows).
	PrevRows []Row
}

// NewTmpFile returns created struct pointer
func NewTmpFile(name string, targetDir string) (*TmpFile, error) {
	exePath, err := os.Executable()
	if err != nil {
		return nil, err
	}
	outPath := filepath.Join(filepath.Dir(exePath), name)

	return &TmpFile{
		Name:      name,
		OutPath:   outPath,
		TargetDir: targetDir,
	}, nil
}

// Create the temp file.
// return the value that has file list information on the search target.
func (f *TmpFile) Create() error {
	file, err := os.OpenFile(f.OutPath, os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	f.Rows, err = f.readRows()
	if err != nil {
		return err
	}

	nlc := getNewLineCode()

	for _, p := range f.Rows {
		_, err := file.WriteString(p.Name + nlc)
		if err != nil {
			return err
		}
	}
	return nil
}

// readRows reads the file list in the direcotory.
func (f TmpFile) readRows() ([]Row, error) {
	files, err := ioutil.ReadDir(f.TargetDir)
	if err != nil {
		return nil, err
	}

	var res []Row
	index := 0
	for _, file := range files {
		if !file.IsDir() {
			row := Row{
				Path: filepath.Join(f.TargetDir, file.Name()),
				Name: file.Name(),
			}
			res = append(res, row)
			index++
		}
	}
	return res, nil
}

// OpenWithEditor opens the this temporary file with editor of specified name.
// After edited it, *TmpFile#Rows is to be edited new rows
// and *TmpFile#PrevRows is to be previous rows.
func (f *TmpFile) OpenWithEditor(editor string) error {
	execCmd := exec.Command(editor, f.OutPath)
	err := execCmd.Run()
	if err != nil {
		return err
	}
	f.PrevRows = f.Rows
	f.Rows, err = f.openRead()
	if err != nil {
		return err
	}
	return nil
}

// openRead opens this temporary file and returns the read Rows.
func (f *TmpFile) openRead() ([]Row, error) {
	var props []Row

	file, err := os.Open(f.OutPath)
	if err != nil {
		return props, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		row := scanner.Text()
		props = append(props, Row{
			Path: filepath.Join(f.TargetDir, row),
			Name: row,
		})
	}

	return props, nil
}

// Remove the this temporary file.
func (f *TmpFile) Remove() error {
	return os.Remove(f.OutPath)
}

// Diff returns difference the previous rows and current rows
func (f *TmpFile) Diff() ([]RowDiff, error) {
	var diffs []RowDiff
	if len(f.PrevRows) != len(f.Rows) {
		return diffs, fmt.Errorf("f.PrevRows and f.Rows are must same length. f.PrevRows length is %d f.Rows length is %d", len(f.PrevRows), len(f.Rows))
	}
	for i := range f.PrevRows {
		if f.PrevRows[i].Name != f.Rows[i].Name {
			diffs = append(diffs, RowDiff{
				PrevRow: f.PrevRows[i],
				CurRow:  f.Rows[i],
			})
		}
	}
	return diffs, nil
}
