package flistfile

// FileRowDiff ファイル行の差分用構造体
type FileRowDiff struct {
	File1 FileRowProp
	File2 FileRowProp
}

// Diff はf1とf2の差分を返却する
func Diff(f1 []FileRowProp, f2 []FileRowProp) []FileRowDiff {
	var diffs []FileRowDiff
	for i := range f1 {
		if f1[i].Name != f2[i].Name {
			diffs = append(diffs, FileRowDiff{
				File1: f1[i],
				File2: f2[i],
			})
		}
	}
	return diffs
}
