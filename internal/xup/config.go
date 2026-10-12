package xup

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/term"
)

type Config struct {
	Endpoint        string `json:"endpoint"`
	Bucket          string `json:"bucket"`
	Region          string `json:"region,omitempty"`
	AccessKeyID     string `json:"access_key_id"`
	SecretAccessKey string `json:"secret_access_key"`
	SessionToken    string `json:"session_token,omitempty"`
	PathStyle       bool   `json:"path_style"`
}

func ConfigPath() (string, error) {
	if v := os.Getenv("XUP_CONFIG"); v != "" {
		return filepath.Abs(v)
	}
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "xup", "config.json"), nil
}
func (c Config) Valid() error {
	if c.Endpoint == "" || c.Bucket == "" || c.AccessKeyID == "" || c.SecretAccessKey == "" {
		return errors.New("endpoint, bucket, access key and secret key are required")
	}
	u, err := url.Parse(c.Endpoint)
	if err != nil || (u.Scheme != "https" && u.Scheme != "http") || u.Host == "" || u.User != nil || u.RawQuery != "" || u.Fragment != "" || (u.Path != "" && u.Path != "/") {
		return errors.New("endpoint must be an http(s) URL without path, credentials, query or fragment")
	}
	if strings.ContainsAny(c.Bucket, "/\\ ") {
		return errors.New("bucket must be a single bucket name")
	}
	return nil
}
func Load() (Config, error) {
	p, err := ConfigPath()
	if err != nil {
		return Config{}, err
	}
	b, err := os.ReadFile(p)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Config{}, fmt.Errorf("not configured; run xup setup: %w", err)
		}
		return Config{}, err
	}
	var c Config
	if err = json.Unmarshal(b, &c); err != nil {
		return c, err
	}
	return c, c.Valid()
}
func Save(c Config) error {
	if err := c.Valid(); err != nil {
		return err
	}
	p, err := ConfigPath()
	if err != nil {
		return err
	}
	if err = os.MkdirAll(filepath.Dir(p), 0700); err != nil {
		return err
	}
	b, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	b = append(b, '\n')
	f, err := os.CreateTemp(filepath.Dir(p), ".config-*")
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())
	defer f.Close()
	if err = f.Chmod(0600); err != nil {
		return err
	}
	if _, err = f.Write(b); err != nil {
		return err
	}
	if err = f.Sync(); err != nil {
		return err
	}
	if err = f.Close(); err != nil {
		return err
	}
	return os.Rename(f.Name(), p)
}
func Prompt(in io.Reader, out io.Writer) (Config, error) {
	var c Config
	read := func(label string, secret bool) (string, error) {
		fmt.Fprint(out, label+": ")
		if secret {
			if f, ok := in.(*os.File); ok && term.IsTerminal(int(f.Fd())) {
				value, err := term.ReadPassword(int(f.Fd()))
				fmt.Fprintln(out)
				return strings.TrimSpace(string(value)), err
			}
		}
		var line string
		_, err := fmt.Fscanln(in, &line)
		if err != nil {
			return "", err
		}
		return strings.TrimSpace(line), nil
	}
	var err error
	if c.Endpoint, err = read("Endpoint (https://...)", false); err != nil {
		return c, err
	}
	if c.Bucket, err = read("Bucket", false); err != nil {
		return c, err
	}
	if c.Region, err = read("Region (enter '-' to omit)", false); err != nil {
		return c, err
	}
	if c.Region == "-" {
		c.Region = ""
	}
	if c.AccessKeyID, err = read("Access key ID", true); err != nil {
		return c, err
	}
	if c.SecretAccessKey, err = read("Secret access key", true); err != nil {
		return c, err
	}
	v, err := read("Path-style addressing? (Y/n)", false)
	if err != nil {
		return c, err
	}
	c.PathStyle = !strings.EqualFold(v, "n")
	return c, c.Valid()
}
func FromEnv() (Config, error) {
	c := Config{Endpoint: os.Getenv("XUP_ENDPOINT"), Bucket: os.Getenv("XUP_BUCKET"), Region: os.Getenv("XUP_REGION"), AccessKeyID: os.Getenv("XUP_ACCESS_KEY_ID"), SecretAccessKey: os.Getenv("XUP_SECRET_ACCESS_KEY"), SessionToken: os.Getenv("XUP_SESSION_TOKEN"), PathStyle: true}
	if v := os.Getenv("XUP_PATH_STYLE"); v != "" {
		if v == "false" || v == "0" {
			c.PathStyle = false
		} else if v != "true" && v != "1" {
			return c, errors.New("XUP_PATH_STYLE must be true/false or 1/0")
		}
	}
	return c, c.Valid()
}
