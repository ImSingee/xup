package xup

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSaveLoad(t *testing.T) {
	t.Setenv("XUP_CONFIG", filepath.Join(t.TempDir(), "config.json"))
	in := Config{Endpoint: "https://s3.example.com", Bucket: "objects", AccessKeyID: "id", SecretAccessKey: "secret", PathStyle: true}
	if err := Save(in); err != nil {
		t.Fatal(err)
	}
	p, _ := ConfigPath()
	stat, err := os.Stat(p)
	if err != nil {
		t.Fatal(err)
	}
	if stat.Mode().Perm() != 0600 {
		t.Fatalf("permissions: %v", stat.Mode().Perm())
	}
	got, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if got != in {
		t.Fatalf("got %+v", got)
	}
}
func TestValidation(t *testing.T) {
	c := Config{Endpoint: "https://s3.example.com/path", Bucket: "test", AccessKeyID: "a", SecretAccessKey: "b"}
	if c.Valid() == nil {
		t.Fatal("accepted endpoint path")
	}
	c.Endpoint = "http://example.com"
	if err := c.Valid(); err != nil {
		t.Fatal(err)
	}
}
func TestHeadless(t *testing.T) {
	t.Setenv("XUP_ENDPOINT", "https://s3.example.com")
	t.Setenv("XUP_BUCKET", "test")
	t.Setenv("XUP_ACCESS_KEY_ID", "a")
	t.Setenv("XUP_SECRET_ACCESS_KEY", "b")
	c, err := FromEnv()
	if err != nil || !c.PathStyle {
		t.Fatalf("got %+v, %v", c, err)
	}
	t.Setenv("XUP_PATH_STYLE", "false")
	c, err = FromEnv()
	if err != nil || c.PathStyle {
		t.Fatalf("got %+v, %v", c, err)
	}
}
func TestSetupOutputDoesNotExposeSecret(t *testing.T) {
	t.Setenv("XUP_CONFIG", filepath.Join(t.TempDir(), "config.json"))
	for k, v := range map[string]string{"XUP_ENDPOINT": "https://s3.example.com", "XUP_BUCKET": "test", "XUP_ACCESS_KEY_ID": "id", "XUP_SECRET_ACCESS_KEY": "super-secret"} {
		t.Setenv(k, v)
	}
	exe := filepath.Join(t.TempDir(), "xup")
	if err := os.WriteFile(exe, []byte("binary"), 0700); err != nil {
		t.Fatal(err)
	}
	var out strings.Builder
	if code := run("xup", []string{"setup", "--headless"}, strings.NewReader(""), &out, &out, func() (string, error) { return exe, nil }); code != 0 {
		t.Fatal(out.String())
	}
	if strings.Contains(out.String(), "super-secret") {
		t.Fatal("secret leaked")
	}
}
