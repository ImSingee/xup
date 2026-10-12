package xup

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestProgramName(t *testing.T) {
	for _, tc := range []struct{ input, want string }{
		{"/usr/local/bin/xup", "xup"},
		{"/usr/local/bin/xdown", "xdown"},
		{`C:\Tools\xdown.exe`, "xdown"},
	} {
		// filepath.Base follows the host OS path conventions.
		if runtime.GOOS != "windows" && tc.input[0] == 'C' {
			continue
		}
		if got := ProgramName(tc.input); got != tc.want {
			t.Errorf("ProgramName(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}

func TestEnsureAlias(t *testing.T) {
	dir := t.TempDir()
	suffix := ""
	if runtime.GOOS == "windows" {
		suffix = ".exe"
	}
	exe := filepath.Join(dir, "xup"+suffix)
	if err := os.WriteFile(exe, []byte("binary"), 0700); err != nil {
		t.Fatal(err)
	}
	alias, err := EnsureAlias(exe)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := EnsureAlias(exe); err != nil {
		t.Fatalf("idempotency: %v", err)
	}
	a, err := os.Stat(alias)
	if err != nil {
		t.Fatal(err)
	}
	b, err := os.Stat(exe)
	if err != nil {
		t.Fatal(err)
	}
	if !os.SameFile(a, b) {
		t.Fatal("alias must refer to xup")
	}
}

func TestEnsureAliasDoesNotOverwrite(t *testing.T) {
	dir := t.TempDir()
	suffix := ""
	if runtime.GOOS == "windows" {
		suffix = ".exe"
	}
	exe := filepath.Join(dir, "xup"+suffix)
	alias := filepath.Join(dir, "xdown"+suffix)
	if err := os.WriteFile(exe, []byte("binary"), 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(alias, []byte("keep"), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := EnsureAlias(exe); err == nil {
		t.Fatal("expected conflict")
	}
	content, err := os.ReadFile(alias)
	if err != nil || string(content) != "keep" {
		t.Fatalf("alias overwritten: %q, %v", content, err)
	}
}
