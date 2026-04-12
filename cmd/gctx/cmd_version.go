package main

import "fmt"

func cmdVersion() int {
	fmt.Printf("gctx version %s %s\n", version, buildDate)
	return 0
}
