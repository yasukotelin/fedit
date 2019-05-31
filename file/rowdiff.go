package file

// RowDiff is struct for difference file row
type RowDiff struct {
	File1 Row
	File2 Row
}

// Diff はf1とf2の差分を返却する
func Diff(f1 []Row, f2 []Row) []RowDiff {
	var diffs []RowDiff
	for i := range f1 {
		if f1[i].Name != f2[i].Name {
			diffs = append(diffs, RowDiff{
				File1: f1[i],
				File2: f2[i],
			})
		}
	}
	return diffs
}
