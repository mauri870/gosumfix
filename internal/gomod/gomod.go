package gomod

import (
	"os/exec"
)

// Tidy runs go mod tidy in the given directory.
func Tidy(dir string) error {
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = dir
	return cmd.Run()
}
