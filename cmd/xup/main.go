package main

import (
	"github.com/ImSingee/xup/internal/xup"
	"os"
)

func main() {
	os.Exit(xup.Run(xup.ProgramName(os.Args[0]), os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}
