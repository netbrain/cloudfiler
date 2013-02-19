package web

import (
	"reflect"
	"testing"
)

type TestDataForm struct {
	Str   string
	Int   int
	Bool  bool
	Byte  []byte
	Strs  []string
	Ints  []int
	Bools []bool
}

type TestDataFormWTag struct {
	Str string `name:"my-cool-STRINGName"`
}

type TestDataFormWCamelCase struct {
	MySimpleString string //should match "mySimpleString"
}

func createValueMap(name string, value ...string) map[string][]string {
	m := make(map[string][]string)
	m[name] = value
	return m
}

func TestNewDataInjectorWithNonPtrArg(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("Expected panic")
		}
	}()

	emptyMap := make(map[string][]string)
	NewDataInjector(emptyMap, TestDataForm{})
}

func TestNewDataInjectorWithNilArg(t *testing.T) {
	fdi := NewDataInjector(nil, &TestDataForm{})
	if fdi != nil {
		t.Fatal("Expected nil in return")
	}
}

func TestInjectStrDataToStruct(t *testing.T) {
	src := createValueMap("str", "testval")
	target := &TestDataForm{}

	fdi := NewDataInjector(src, target)
	fdi.Inject()

	if target.Str != "testval" {
		t.Fatal("Data was not injected")
	}
}

func TestInjectIntDataToStruct(t *testing.T) {
	src := createValueMap("int", "10")
	target := &TestDataForm{}

	fdi := NewDataInjector(src, target)
	fdi.Inject()

	if target.Int != 10 {
		t.Fatal("Data was not injected")
	}
}

func TestInjectBoolDataToStruct(t *testing.T) {
	src := createValueMap("bool", "true")
	target := &TestDataForm{}

	fdi := NewDataInjector(src, target)
	fdi.Inject()

	if !target.Bool {
		t.Fatal("Data was not injected")
	}
}

func TestInjectByteDataToStruct(t *testing.T) {
	value := "some-string-to-bytearray"
	src := createValueMap("byte", value)
	target := &TestDataForm{}

	fdi := NewDataInjector(src, target)
	fdi.Inject()

	if len(target.Byte) != len(value) {
		t.Fatal("Data was not injected")
	}
}

func TestInjectSeveralStrDataToStruct(t *testing.T) {
	src := createValueMap("strs", "1", "abc")
	target := &TestDataForm{}

	fdi := NewDataInjector(src, target)
	fdi.Inject()

	if len(target.Strs) != 2 {
		t.Fatal("Data was not injected")
	}

	if !reflect.DeepEqual(
		target.Strs, []string{"1", "abc"}) {
		t.Fatalf("Incorrect data injected, found: %#v", target.Strs)
	}
}

func TestInjectSeveralIntDataToStruct(t *testing.T) {
	src := createValueMap("ints", "1", "2", "3")
	target := &TestDataForm{}

	fdi := NewDataInjector(src, target)
	fdi.Inject()

	if len(target.Ints) != 3 {
		t.Logf("taget.Ints: %#v", target.Ints)
		t.Fatal("Correct data was not injected")
	}

	if !reflect.DeepEqual(
		target.Ints, []int{1, 2, 3}) {
		t.Fatalf("Incorrect data injected, found: %#v", target.Ints)
	}
}

func TestInjectSeveralBoolDataToStruct(t *testing.T) {
	src := createValueMap("bools", "0", "1", "T", "f", "True", "false")
	target := &TestDataForm{}

	fdi := NewDataInjector(src, target)
	fdi.Inject()

	if len(target.Bools) != 6 {
		t.Logf("taget.Bools: %#v", target.Bools)
		t.Fatal("Correct data was not injected")
	}

	if !reflect.DeepEqual(
		target.Bools, []bool{false, true, true, false, true, false}) {
		t.Fatalf("Incorrect data injected, found: %#v", target.Bools)
	}
}

func TestInjectStrDataToStructWTag(t *testing.T) {
	src := createValueMap("my-cool-STRINGName", "testval")
	target := &TestDataFormWTag{}

	fdi := NewDataInjector(src, target)
	fdi.Inject()

	if target.Str != "testval" {
		t.Fatalf("Data was not injected")
	}
}

func TestInjectStrDataToStructWCamelCaseNaming(t *testing.T) {
	src := createValueMap("mySimpleString", "testval")
	target := &TestDataFormWCamelCase{}

	fdi := NewDataInjector(src, target)
	fdi.Inject()

	if target.MySimpleString != "testval" {
		t.Fatalf("Data was not injected")
	}
}
