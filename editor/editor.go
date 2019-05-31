package editor

import "runtime"

// GetDefaultEditor returns default editor name on the running OS.
func GetDefaultEditor() string {
	if runtime.GOOS == "windows" {
		return "notepad"
	}
	return "vim"
}
