package html_from_yaml

import (
	"testing"
)

var testyaml string = `
html:
  attrs: []
  content:
    - head:
        attrs: []
        content:
          - meta:
              attrs:
                - name: 'xyz'
              content: []
          - title:
              attrs: []
              content:
                - text: 'xyz'
    - body:
        attrs: []
        content:
          - div:
              attrs:
                - name: 'id'
                  value: 'bob'
              content:
                - text: 'xyz'
          - div:
              attrs:
                - name: 'id'
                  value: 'dud'
              content:
                - img:
                    attrs:
                      - name: 'src'
                        value: 'xyz'
                    content: []
`

/*
<html>
	<head>
		<meta xyz />
		<title>xyz</title>
	</head>
	<body>
		<div id="bob">xyz</div>
		<div id="dud"><img src="xyz" /></div>
	</body>
</html>
*/

func TestHtmlFromYaml(t *testing.T) {
	t.Logf(testyaml)
	t.Logf(UsHtmlFromYamlUs(testyaml))
}
