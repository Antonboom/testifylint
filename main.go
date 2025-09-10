package main

import (
	"fmt"
	"os"
	"runtime/debug"
	"slices"

	"golang.org/x/tools/go/analysis/singlechecker"

	"github.com/Antonboom/testifylint/analyzer"
)

// Populated by goreleaser during build.
var (
	version = "unknown"
	commit  = "unknown"
	date    = "unknown"
)

func main() {
	for _, arg := range os.Args[1:] {
		if slices.Contains([]string{"--version", "-version", "-V", "--V"}, arg) {
			printVersion()
			return
		}
	}

	singlechecker.Main(analyzer.New())
}

func printVersion() {
	goVersion := "go???"
	if buildInfo, ok := debug.ReadBuildInfo(); ok {
		goVersion = buildInfo.GoVersion
	}

	fmt.Printf("version %s built with %s from %s on %s\n", version, goVersion, commit, date)
}
