package xup

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// ProgramName decides the command from the name used to invoke this binary.
func ProgramName(arg0 string) string {
	name := strings.ToLower(filepath.Base(arg0))
	name = strings.TrimSuffix(name, ".exe")
	return name
}

// EnsureAlias creates an xdown alias alongside xup without replacing unrelated files.
// On Windows, a hard link is used if symlinks are unavailable to the current user.
func EnsureAlias(executable string) (string, error) {
	executable, err := filepath.Abs(executable)
	if err != nil {
		return "", err
	}
	suffix := ""
	if runtime.GOOS == "windows" {
		suffix = ".exe"
	}
	alias := filepath.Join(filepath.Dir(executable), "xdown"+suffix)
	source, err := os.Stat(executable)
	if err != nil {
		return alias, err
	}
	if !source.Mode().IsRegular() {
		return alias, fmt.Errorf("executable is not a regular file: %s", executable)
	}
	if target, err := os.Stat(alias); err == nil {
		if os.SameFile(source, target) {
			return alias, nil
		}
		return alias, fmt.Errorf("%s already exists and points to a different file; leaving it untouched", alias)
	} else if !os.IsNotExist(err) {
		return alias, fmt.Errorf("cannot inspect %s: %w", alias, err)
	}
	// A dangling symlink exists even though Stat reported it as missing.
	if _, err := os.Lstat(alias); err == nil {
		return alias, fmt.Errorf("%s already exists and points to a different file; leaving it untouched", alias)
	} else if !os.IsNotExist(err) {
		return alias, fmt.Errorf("cannot inspect %s: %w", alias, err)
	}
	if err := os.Symlink(filepath.Base(executable), alias); err != nil {
		if runtime.GOOS == "windows" {
			if linkErr := os.Link(executable, alias); linkErr == nil {
				return alias, nil
			}
		}
		return alias, fmt.Errorf("cannot create xdown alias next to %s: %w (install xup in a writable directory or create the alias manually)", executable, err)
	}
	return alias, nil
}
