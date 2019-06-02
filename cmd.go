package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yasukotelin/fedit/editor"
	"github.com/yasukotelin/fedit/file"
)

var (
	specifiedEditor string

	rootCmd = &cobra.Command{
		Use:   "fedit",
		Short: "rename all files in derectory.",
		Long:  "fedit is the tool to rename all files in directory",
		Run: func(cmd *cobra.Command, args []string) {
			if err := run(args); err != nil {
				fmt.Printf("Error: %s\n", err)
				os.Exit(ExitCodeError)
			}
			os.Exit(ExitCodeError)
		},
	}
)

func init() {
	rootCmd.Version = version
	rootCmd.PersistentFlags().StringVarP(&specifiedEditor, "editor", "e", editor.GetDefaultEditor(), "specify the editor to open. ")
}

// run does the following flow.
//
// 1. create the temp file of file list.
// 2. open it with the specified editor and user can edit it.
// 3. read the edited it and get edited file name list.
// 4. do rename them.
// 5. delete the temp file.
func run(args []string) error {
	dirPath, err := getDirPath(args)
	if err != nil {
		return err
	}
	tmpFile, err := file.NewTmpFile(tmpFileName, dirPath)
	if err != nil {
		return err
	}

	if err = tmpFile.Create(); err != nil {
		return err
	}
	defer tmpFile.Remove()

	if err := tmpFile.OpenWithEditor(specifiedEditor); err != nil {
		return err
	}

	if tmpFile.IsDeletedRows() {
		return errors.New("Deleted rows error")
	}

	if tmpFile.IsAddedRows() {
		return errors.New("Added new rows error")
	}

	diffs, err := tmpFile.Diff()
	if err != nil {
		return err
	}

	switch {
	case len(diffs) == 0:
		return errors.New("no changed the file name")
	case len(diffs) > 0:
		if err = doRenameWithConfirm(diffs); err != nil {
			return err
		}
	}

	if err := tmpFile.Remove(); err != nil {
		return err
	}

	return nil
}

func getDirPath(args []string) (string, error) {
	if len(args) == 0 {
		return "", errors.New("directry path required")
	}
	return args[0], nil
}

func printDiff(diffs []file.RowDiff) {
	fmt.Println()
	for _, d := range diffs {
		fmt.Printf("%s ---> %s\n", d.PrevRow.Path, d.CurRow.Path)
	}
	fmt.Println()
}

func askToApplyRename() (bool, error) {
	var s string
	for {
		fmt.Print("apply to rename [y/n]? ")
		_, err := fmt.Scanln(&s)
		if err != nil {
			return false, err
		}
		switch s {
		case "y", "Y":
			return true, nil
		case "n", "N":
			return false, nil
		}
	}
}

func doRenameWithConfirm(diffs []file.RowDiff) error {
	printDiff(diffs)

	ok, err := askToApplyRename()
	if err != nil {
		return err
	}
	if ok {
		err = rename(diffs)
		if err != nil {
			return err
		}
	}
	return nil
}

func rename(diffs []file.RowDiff) error {
	for _, d := range diffs {
		if err := os.Rename(d.PrevRow.Path, d.CurRow.Path); err != nil {
			return err
		}
	}
	return nil
}
