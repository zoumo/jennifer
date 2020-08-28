package jen

import (
	"fmt"
	"testing"
)

func TestGuessAlias(t *testing.T) {

	data := map[string]string{
		"A":             "a",
		"a":             "a",
		"a$":            "a",
		"a/b":           "b",
		"a/b/c":         "c",
		"a/b/c-d":       "cd",
		"a/b/c-d/":      "cd",
		"a.b":           "ab",
		"a/b.c":         "bc",
		"a/b-c.d":       "bcd",
		"a/bb-ccc.dddd": "bbcccdddd",
		"a/foo-go":      "foogo",
		"123a":          "a",
		"a/321a.b":      "ab",
		"a/123":         "pkg",
	}
	for path, expected := range data {
		if guessAlias(path) != expected {
			fmt.Printf("guessAlias test failed %s should return %s but got %s\n", path, expected, guessAlias(path))
			t.Fail()
		}
	}
}

func TestValidAlias(t *testing.T) {
	data := map[string]bool{
		"a":   true,  // ok
		"b":   false, // already registered
		"go":  false, // keyword
		"int": false, // predeclared
		"err": false, // common name
	}
	f := NewFile("test")
	f.register("b")
	for alias, expected := range data {
		if f.isValidAlias(alias) != expected {
			fmt.Printf("isValidAlias test failed %s should return %t but got %t\n", alias, expected, f.isValidAlias(alias))
			t.Fail()
		}
	}
}

func TestFile_register(t *testing.T) {
	data := newOrderedMap()
	data.Add("v1", "v1")
	data.Add("meta.v1", "metav1")
	data.Add("metav1", "metav1x")
	data.Add("meta-v1", "metav1xx")
	data.Add("meta/v1", "metav1xxx")
	data.Add("github.com/xxx/foo.abc", "fooabc")
	data.Add("github.com/xxx/fooabc", "xxxfooabc")
	data.Add("github.com/xxx/fooa.bc", "githubcomxxxfooabc")
	data.Add("aaa/bbb/123ccc", "ccc")
	data.Add("aaa/123bbb/123ccc", "bbbccc")

	f := NewFile("test")

	data.Range(func(path, alias string) bool {
		got := f.register(path)
		if got != alias {
			fmt.Printf("register test failed %s should return %s but got %s\n", path, alias, got)
			t.Fail()
		}
		return true
	})
}

type orderedMap struct {
	order []string
	data  map[string]string
}

func newOrderedMap() *orderedMap {
	return &orderedMap{
		order: make([]string, 0),
		data:  make(map[string]string),
	}
}

func (m *orderedMap) Add(k, v string) {
	_, ok := m.data[k]
	if ok {
		return
	}

	m.data[k] = v
	m.order = append(m.order, k)
}

func (m *orderedMap) Range(visit func(k, v string) bool) {
	for _, k := range m.order {
		v := m.data[k]
		if !visit(k, v) {
			break
		}
	}
}
