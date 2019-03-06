package flistfile

// Row is struct that has infomation for one line of tmp file
type Row struct {
	Path string
	Name string
}

// IsDupl returns whether the rows are duplicated.
func IsDupl(rows *[]Row) bool {
	m := make(map[string]bool)
	for _, row := range *rows {
		if m[row.Path] {
			return true
		}
		m[row.Path] = true
	}
	return false
}
