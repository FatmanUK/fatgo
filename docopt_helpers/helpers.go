package docopt_helpers

import (
	"strings"
	"github.com/docopt/docopt-go"
)

type HelpOption struct {
	Value string
	Desc string
	Default string
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

func MakeUses(usages []string, options map[string]HelpOption) ([]string, map[string]string) {
	opts := make(map[string]string)

	for k1, u1 := range usages {
		for k2, o2 := range options {
			k := o2.MakeOption(k2)
			u1 = strings.Replace(u1, k2, k, -1)
		}
		usages[k1] = u1
	}

	for k, v := range options {
		t := v.MakeOption(k)
		opts[t] = v.Desc
		if v.Default != "" {
			opts[t] += " " + MakeDefault(v.Default)
		}
	}

	return usages, opts
}

