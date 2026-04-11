package gitcmd

import "strings"

// FirstGitSubcommand parses git's global options and returns the first
// non-flag argument (the subcommand) and the slice of arguments after it.
//
// Recognised global options that consume a following value:
//
//	-c <key=value>, -C <path>, --git-dir <path>, --work-tree <path>, --namespace <ns>
//
// Options in --key=value form and bare single-dash flags are consumed
// individually. Unknown global flags with space-separated arguments may cause
// the argument to be misidentified as the subcommand; this is a known
// limitation documented here rather than silently ignored.
func FirstGitSubcommand(args []string) (sub string, restAfterSub []string) {
	i := 0
	for i < len(args) {
		a := args[i]
		switch {
		case a == "-c" && i+1 < len(args):
			i += 2
		case a == "-C" || a == "--git-dir" || a == "--work-tree" || a == "--namespace":
			if i+1 < len(args) {
				i += 2
			} else {
				i++
			}
		case strings.HasPrefix(a, "--git-dir=") || strings.HasPrefix(a, "--work-tree=") ||
			strings.HasPrefix(a, "--namespace="):
			i++
		case strings.HasPrefix(a, "-") && a != "-":
			i++
		default:
			return a, args[i+1:]
		}
	}
	return "", nil
}
