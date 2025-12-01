package doh

import (
	"bytes"
	"text/template"
	"github.com/docopt/docopt-go"
)

var NoExitParser = &docopt.Parser{
	HelpHandler:   docopt.PrintHelpOnly,
	OptionsFirst:  false,
	SkipHelpFlags: false,
}

type DocOptTemplate struct {
	t *template.Template
	wr bytes.Buffer
}

func (re DocOptTemplate) Init() DocOptTemplate {
	re.t = template.New("DocOptTemplate")
	return re
}

func (re DocOptTemplate) Execute (
		name string,
		values interface{}) (string, error) {
	var err error
	re.t, err = re.t.Parse(name)
	if err != nil {
		return "", nil
	}
	err = re.t.Execute(&re.wr, values)
	if err != nil {
		return "", nil
	}
	return string(re.wr.Bytes()), nil
}

func RunDocopt(
		docoptAppVer string,
		docoptString string,
		dotv interface{}) (map[string]interface{}, error) {
	dot := DocOptTemplate{}.Init()
	dav, err := dot.Execute(docoptAppVer, dotv)
	if err != nil {
		return nil, err
	}
	ds, err := dot.Execute(docoptString, dotv)
	if err != nil {
		return nil, err
	}
	return NoExitParser.ParseArgs(ds, nil, dav)
}
