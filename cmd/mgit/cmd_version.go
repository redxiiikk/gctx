package main

import "fmt"

func cmdVersion() int {
	fmt.Printf("mgit version %s %s\n", version, buildDate)
	return 0
}
