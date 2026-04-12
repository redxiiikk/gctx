package main

import "fmt"

var (
	version = "dev-version"
	commit  = "none"
	date    = "unknown"
)

func cmdVersion() int {
	fmt.Printf("gctx version %s, commit %s, built at %s\n", version, commit, date)
	return 0
}
