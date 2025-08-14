package docopt_helpers

import (
	"os"
	"path"
	"strings"
)

func BuildUsageString(uses []string, opts map[string]string) string {
	base := path.Base(os.Args[0])

	// Add some defaults
	uses = append(uses, "-h | --help")
	uses = append(uses, "-v | --version")
	opts["-h --help"] = "Show this screen"
	opts["-v --version"] = "Show version"
	var rv string
	// Start with app name. Arg0 basename followed by dot.
	rv = base + ".\n\nUsage:\n"
	// Add usage strings.
	for _, use := range uses {
		rv += "  " + base + " " + use + "\n"
	}
	// Add option strings.
	rv += "\nOptions:\n"
	max_len := 0
	for first, _ := range opts {
		cur_len := len(first)
		if cur_len > max_len {
			max_len = cur_len
		}
	}
	for first, second := range opts {
		pad := max_len - len(first)
		rv += "  " + first
		rv += strings.Repeat(" ", pad + 2)
		rv += second + ".\n"
	}
	return rv
}
