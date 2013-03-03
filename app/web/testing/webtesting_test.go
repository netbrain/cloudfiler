package testing

import (
	"testing"
)

func TestCreateReqContext(t *testing.T) {
	ctx, _ := CreateReqContext("POST", "/test/path", nil)
	if ctx.Method() != "POST" {
		t.Fatal("Expected POST as method")
	}

	p := ctx.Request.URL.Path
	if p != "/test/path" {
		t.Fatal("Expected /test/path")
	}
}
