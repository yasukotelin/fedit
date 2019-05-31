package file

import "testing"

func TestIsDupl(t *testing.T) {
	// 重複していない
	rows := []Row{
		Row{Path: "./diff.go"},
		Row{Path: "./Row.go"},
		Row{Path: "./flistfile.go"},
	}

	act := IsDupl(&rows)
	exec := false

	if act != exec {
		t.Errorf("Error, it returns true even though it's not duplicated")
	}

	// 重複している
	rows = []Row{
		Row{Path: "./diff.go"},
		Row{Path: "./Row.go"},
		Row{Path: "./flistfile.go"},
		Row{Path: "./diff.go"},
	}

	act = IsDupl(&rows)
	exec = true

	if act != exec {
		t.Errorf("Error, it returns true even though it's not duplicated")
	}
}
