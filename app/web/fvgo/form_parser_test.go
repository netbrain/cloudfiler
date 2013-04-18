package fvgo

import (
	"code.google.com/p/go.net/html"
	"reflect"
	"strings"
	"testing"
)

func TestGetAttrsMapFromNode(t *testing.T) {
	htmlText := `<form action="/some/action"></form>`
	node, _ := html.Parse(strings.NewReader(htmlText))
	node = findChildrenNodesByTag(node, "form")[0]

	result := getAttrsMapFromNode(node)
	expected := map[string]string{
		"action": "/some/action",
	}

	if !reflect.DeepEqual(result, expected) {
		t.Logf("%#v != %#v", result, expected)
		t.Fatal("Failed on equality")
	}
}

func TestParseField(t *testing.T) {
	htmlText := `<input type="text" name="test" />`
	node, _ := html.Parse(strings.NewReader(htmlText))
	node = findChildrenNodesByTag(node, "input")[0]

	result := parseField(node)

	if result.tag != "input" {
		t.Fatalf("Expected 'input' but got: %s", result.tag)
	}

	expected := getAttrsMapFromNode(node)
	if !reflect.DeepEqual(result.attrs, expected) {
		t.Logf("%#v != %#v", result, expected)
		t.Fatal("Failed on equality")
	}
}

func TestParseForm(t *testing.T) {
	htmlText := `<form action="/some/action"></form>`
	node, _ := html.Parse(strings.NewReader(htmlText))
	node = findChildrenNodesByTag(node, "form")[0]

	result := parseForm(node)

	expected := getAttrsMapFromNode(node)
	if !reflect.DeepEqual(result.attrs, expected) {
		t.Logf("%#v != %#v", result, expected)
		t.Fatal("Failed on equality")
	}
}

func TestParse(t *testing.T) {
	htmlText := `
	  <form action="/some/action">
        <input type="text" name="test" />
	  </form>
	`
	root, _ := html.Parse(strings.NewReader(htmlText))
	result := parse(root)

	if len(result) != 1 {
		t.Fatalf("Expected one form, got: %v", len(result))
	}

	if len(result[0].fields) != 1 {
		t.Fatalf("Expected one field, got: %v", len(result))
	}
}

func TestParseWithMultipleFields(t *testing.T) {
	htmlText := `
	  <form action="/some/action">
        <input type="text" name="test" />
        <input type="text" name="test2" />
		<input type="text" name="test3" />        
	  </form>
	`
	root, _ := html.Parse(strings.NewReader(htmlText))
	result := parse(root)

	if len(result) != 1 {
		t.Fatalf("Expected one form, got: %v", len(result))
	}

	if len(result[0].fields) != 3 {
		t.Fatalf("Expected three field, got: %v", len(result))
	}
}
