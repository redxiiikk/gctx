package repos

import "fmt"

func CmdRepos(args []string) int {
	switch args[0] {
	case "help":
		fmt.Println("This is repos command")
		return 0
	}
	return 1
}
