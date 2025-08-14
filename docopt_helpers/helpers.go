package docopt_helpers

import (
	"github.com/docopt/docopt-go"
)

type HelpOption struct {
	Value string
	Desc string
}

var NoExitParser = &docopt.Parser{
	HelpHandler:   docopt.PrintHelpOnly,
	OptionsFirst:  false,
	SkipHelpFlags: false,
}

func IsSet(arr map[string]interface{}, key string) bool {
	_, ok := arr[key]
	return ok
}

func MakeDefault(value string) string {
	return "[default: " + value + "]"
}

func (re HelpOption) MakeOption(flag string) string {
	return flag + "=" + re.Value
}

