// Package version provides ways to get the version of the
// program.
package version

import "os/exec"

// Returns the commit hash of the repository.
func CommitHash() string {
	cmd := exec.Command("git", []string{"rev-parse", "HEAD"}...)

	b, err := cmd.Output()
	if err != nil {
		return "unknown"
	}
	return string(b)
}
