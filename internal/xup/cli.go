package xup

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

const usage = `Usage:
  xup setup [--headless]      Configure S3 and create an xdown alias beside xup
  xup FILE [KEY]             Upload a local file (default key: filename)
  xup rm KEY                 Delete a remote object
  xup ls [PREFIX]            List remote objects
  xdown KEY [FILE]           Download (default file: key basename; never overwrites)

Environment for headless setup: XUP_ENDPOINT, XUP_BUCKET, XUP_ACCESS_KEY_ID,
XUP_SECRET_ACCESS_KEY; optional XUP_REGION, XUP_SESSION_TOKEN, XUP_PATH_STYLE.
XUP_CONFIG overrides the config file path. Use HTTPS unless the endpoint is trusted.
`

func Run(program string, args []string, in io.Reader, out, errout io.Writer) int {
	return run(program, args, in, out, errout, os.Executable)
}

func run(program string, args []string, in io.Reader, out, errout io.Writer, executable func() (string, error)) int {
	if program != "xup" && program != "xdown" {
		fmt.Fprintf(errout, "unsupported command name %q: run as xup or xdown\n", program)
		return 2
	}
	if len(args) == 0 {
		fmt.Fprint(errout, usage)
		return 2
	}
	if args[0] == "help" || args[0] == "--help" || args[0] == "-h" {
		fmt.Fprint(out, usage)
		return 0
	}
	if program == "xup" && args[0] == "setup" {
		flags := flag.NewFlagSet("setup", flag.ContinueOnError)
		flags.SetOutput(errout)
		headless := flags.Bool("headless", false, "read configuration from environment variables")
		if err := flags.Parse(args[1:]); err != nil {
			return 2
		}
		if flags.NArg() != 0 {
			fmt.Fprintln(errout, "setup takes no positional arguments")
			return 2
		}
		var c Config
		var err error
		if *headless {
			c, err = FromEnv()
		} else {
			if f, ok := in.(*os.File); ok {
				st, e := f.Stat()
				if e != nil || st.Mode()&os.ModeCharDevice == 0 {
					fmt.Fprintln(errout, "non-interactive stdin: use xup setup --headless")
					return 2
				}
			}
			c, err = Prompt(in, out)
		}
		if err == nil {
			err = Save(c)
		}
		if err != nil {
			fmt.Fprintln(errout, "setup:", err)
			return 1
		}
		p, _ := ConfigPath()
		fmt.Fprintln(out, "Configuration saved to", p)
		exe, e := executable()
		if e != nil {
			fmt.Fprintln(errout, "setup: configuration saved, but cannot locate executable for xdown alias:", e)
			return 1
		}
		alias, e := EnsureAlias(exe)
		if e != nil {
			fmt.Fprintln(errout, "setup: configuration saved, but cannot create xdown alias:", e)
			return 1
		}
		fmt.Fprintln(out, "xdown alias ready at", alias)
		return 0
	}
	c, err := Load()
	if err != nil {
		fmt.Fprintln(errout, err)
		return 1
	}
	store, err := Connect(c)
	if err != nil {
		fmt.Fprintln(errout, err)
		return 1
	}
	ctx := context.Background()
	switch {
	case program == "xdown":
		if len(args) < 1 || len(args) > 2 || strings.HasPrefix(args[0], "-") {
			fmt.Fprint(errout, usage)
			return 2
		}
		key := args[0]
		local := DefaultLocal(key)
		if len(args) == 2 {
			local = args[1]
		}
		if local == "." || local == "" {
			fmt.Fprintln(errout, "invalid destination filename")
			return 2
		}
		n, e := store.Download(ctx, key, local)
		err = e
		if err == nil {
			fmt.Fprintf(out, "Downloaded %s -> %s (%d bytes)\n", key, local, n)
		}
	case program == "xup" && args[0] == "rm":
		if len(args) != 2 || args[1] == "" {
			fmt.Fprint(errout, usage)
			return 2
		}
		err = store.Remove(ctx, args[1])
		if err == nil {
			fmt.Fprintln(out, "Deleted", args[1])
		}
	case program == "xup" && args[0] == "ls":
		if len(args) > 2 {
			fmt.Fprint(errout, usage)
			return 2
		}
		prefix := ""
		if len(args) == 2 {
			prefix = args[1]
		}
		err = store.List(ctx, prefix, out)
	case program == "xup":
		if len(args) > 2 || args[0] == "" || strings.HasPrefix(args[0], "-") {
			fmt.Fprint(errout, usage)
			return 2
		}
		key := DefaultKey(args[0])
		if len(args) == 2 {
			key = args[1]
		}
		if key == "" {
			fmt.Fprintln(errout, "empty object key")
			return 2
		}
		n, e := store.Upload(ctx, args[0], key)
		err = e
		if err == nil {
			fmt.Fprintf(out, "Uploaded %s -> %s (%d bytes)\n", args[0], key, n)
		}
	default:
		err = errors.New("unknown executable")
	}
	if err != nil {
		fmt.Fprintln(errout, "error:", err)
		return 1
	}
	return 0
}
