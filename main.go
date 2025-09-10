package main

import (
	"flag"
	"fmt"
	"runtime/debug"

	"golang.org/x/tools/go/analysis/singlechecker"

	"github.com/Antonboom/testifylint/analyzer"
)

// Populated by goreleaser during build.
const (
	version = "unknown"
	commit
	date
)

func main() {
	showVersion := flag.Bool("V", false, "print version and exit")
	flag.Parse()

	if *showVersion {
		goVersion := "go???"
		if buildInfo, ok := debug.ReadBuildInfo(); ok {
			goVersion = buildInfo.GoVersion
		}

		fmt.Printf("version %s built with %s from %s on %s\n", version, goVersion, commit, date)
		return
	}

	singlechecker.Main(analyzer.New())
}
