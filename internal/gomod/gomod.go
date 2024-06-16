package gomod

import (
	"fmt"
	"os/exec"
)

// Tidy runs go mod tidy in the given directory.
func Tidy(dir string) error {
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to run go mod tidy: %v:\n%s", err, out)
	}
	return nil
}
