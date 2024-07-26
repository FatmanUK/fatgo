package html_from_yaml

import (
	"fmt"
//	"bytes"
//	"os"
//	"path"
//	"strings"
//	"golang.org/x/net/html"
//	"net/html"
	"gopkg.in/yaml.v3"
)

/*func getNode(parent map[string]interface{}) map[string]interface{} {


	return data["html"].(map[string]interface{})

}*/

func getTopNode(d map[string]interface{}, h string) map[string]interface{} {
	return d[h].(map[string]interface{})
}

func getNode(d interface{}, h string) map[string]interface{} {
	return d.(map[string]interface{})[h].(map[string]interface{})
}

func getContent(d map[string]interface{}) []interface{} {
	return d["content"].([]interface{})
}

type Attr struct {
	Name string
	Value string
}

func getAttrs(d map[string]interface{}) []Attr {
	attrs := []Attr{}
	for _, v := range d["attrs"].([]interface{}) {
		w := v.(map[string]interface{})
		a := &Attr{Name: w["name"].(string), Value: w["value"].(string)}
		attrs = append(attrs, *a)
	}
	return attrs
}

func UsHtmlFromYamlUs(usYaml string) string {
	// Map to store the parsed YAML data
	var data map[string]interface{}

	// Unmarshal the YAML string into the map
	err := yaml.Unmarshal([]byte(usYaml), &data)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println()

	html := getTopNode(data, "html")
	//fmt.Println("html:", html)
	//html_attrs := getAttrs(html)
	//fmt.Println("attrs:", html_attrs)
	html_content := getContent(html)
//	head := html_content[0].(map[string]interface{})["head"]
	head := getNode(html_content[0], "head")
	body := getNode(html_content[1], "body")
	fmt.Println("head:", head)
	fmt.Println("body:", body)
//	fmt.Println("City:", data["address"].(map[string]interface{})["city"])

	return ""
}
