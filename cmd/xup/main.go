package main

import (
	"github.com/ImSingee/xup/internal/xup"
	"os"
)

func main() { os.Exit(xup.Run("xup", os.Args[1:], os.Stdin, os.Stdout, os.Stderr)) }
