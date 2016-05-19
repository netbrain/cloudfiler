package fvgo

import (
	"golang.org/x/net/html"
	"os"
)

func Parse(filePath string) ([]*Form, error) {
	forms := make([]*Form, 0)
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	node, _ := html.Parse(file)
	for _, form := range parse(node) {
		forms = append(forms, form)
	}
	return forms, nil
}

func parse(root *html.Node) []*Form {
	forms := make([]*Form, 0)
	formNodes := findChildrenNodesByTag(root, "form")
	for _, formNode := range formNodes {
		form := parseForm(formNode)
		forms = append(forms, form)

		fieldNodes := findChildrenNodesByTag(formNode, "input")
		for _, fieldNode := range fieldNodes {
			field := parseField(fieldNode)
			form.addField(field)
		}
	}
	return forms
}

func findChildrenNodesByTag(n *html.Node, tag string) []*html.Node {
	nodes := make([]*html.Node, 0)
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		for _, nc := range findChildrenNodesByTag(c, tag) {
			nodes = append(nodes, nc)
		}
		if !(c.Type == html.ElementNode && c.Data == tag) {
			continue
		}
		nodes = append(nodes, c)
	}
	return nodes
}

func parseForm(formNode *html.Node) *Form {
	return &Form{
		attrs: getAttrsMapFromNode(formNode),
	}
}

func parseField(fieldNode *html.Node) *field {
	return NewField(fieldNode.Data, getAttrsMapFromNode(fieldNode))
}

func getAttrsMapFromNode(n *html.Node) map[string]string {
	attrs := make(map[string]string)
	for _, a := range n.Attr {
		attrs[a.Key] = a.Val
	}
	return attrs
}
