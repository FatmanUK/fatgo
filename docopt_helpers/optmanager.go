package docopt_helpers

import (
	"fmt"
	"path"
	"strings"
)

type Opt struct {
	small rune
	large string
	desc string
}

type OptPattern struct {
	pattern string
	refs []Opt
}

func (re OptPattern) Output() string {
	var rv string = re.pattern
	opts := []interface{}{}
	for _, opt := range re.refs {
		str := fmt.Sprintf("-%s | --%s", string(opt.small), opt.large)
		opts = append(opts, str)
	}
	rv = fmt.Sprintf(rv, opts...)
	return rv
}

type OptManager struct {
	Opts map[rune]Opt
	Patterns []OptPattern
}

func (re OptManager) AddOpt(small rune, large string, desc string) OptManager {
	re.Opts[small] = Opt{small, large, desc}
	return re
}

func (re OptManager) AddPattern(pattern string, refs []Opt) OptManager {
	re.Patterns = append(re.Patterns, OptPattern{pattern, refs})
	return re
}

func (re OptManager) GetBool(d docopt.Opts, r rune) bool {
	small := "-" + string(r)
	large := "--" + re.Opts[r].large
	chk1, err1 := d.Bool(small)
	if err != nil {
		panic("wtf1")
	}
	chk2, err2 := d.Bool(large)
	if err != nil {
		panic("wtf2")
	}
	return chk1 || chk2
}

func (re OptManager) LongestKey() string {
	lk := ""
	for _, opt := range re.Opts {
		if len(opt.large) > len(lk) {
			lk = opt.large
		}
	}
	return lk
}

func (re OptManager) Usage(args0 string) string {
	args0 = path.Base(args0)
	var rv string = fmt.Sprintf("%s.\n\nUsage:\n", args0)
	for _, pattern := range re.Patterns {
		use := pattern.Output()
		rv += fmt.Sprintf("  %s %s\n", args0, use)
	}
	rv += "\nOptions:\n"
	max_len := len(re.LongestKey())
	for _, opt := range re.Opts {
		pad := max_len - len(opt.large)
		rv += fmt.Sprintf(
			"  -%s --%s %s %s\n",
			string(opt.small),
			opt.large,
			strings.Repeat(" ", pad),
			opt.desc,
		)
	}
	return rv
}
