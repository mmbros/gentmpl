package lib

import (
	"fmt"
	"testing"
)

func ExampleresolveIncludes() {
	var (
		mapping = map[string][]string{
			"inc": {"header", "footer"},
			"A":   {"contentA", "inc"},
			"B":   {"contentB", "inc"},
		}
		names = []string{"A", "B"}
	)
	res, _ := ResolveIncludes(mapping, names)
	fmt.Println(res["A"])
	fmt.Println(res["B"])
	// Output:
	// [contentA header footer]
	// [contentB header footer]
}

func checkEqual(a, b []string) bool {

	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func TestResolveIncludes(t *testing.T) {
	var templates = map[string][]string{
		"tmplA": {"A1", "A2", "A3"},
		"incB":  {"B1", "B2"},
		"incC":  {"incB", "C1"},
		"tmplD": {"incB", "D1"},
		"tmplE": {"incC", "E1", "E2"},
		"tmplF": {"incC", "incB", "F1"},
	}
	var astr = []string{"incC", "tmplA", "tmplD", "tmplE", "tmplF"}

	var expected = map[string][]string{
		"incC":  {"B1", "B2", "C1"},
		"tmplA": {"A1", "A2", "A3"},
		"tmplD": {"B1", "B2", "D1"},
		"tmplE": {"B1", "B2", "C1", "E1", "E2"},
		"tmplF": {"B1", "B2", "C1", "B1", "B2", "F1"},
	}
	res, _ := ResolveIncludes(templates, astr)
	for _, tmpl := range astr {
		if !checkEqual(res[tmpl], expected[tmpl]) {
			t.Errorf("resolveIncludes(%q): expected %v, actual %v", tmpl, expected[tmpl], res[tmpl])
		}
	}

}

func TestResolveIncludesLoop(t *testing.T) {
	var templates = map[string][]string{
		"incA": {"A1", "incB"},
		"incB": {"incA", "B1"},
	}
	var astr = []string{"incA"}

	_, err := ResolveIncludes(templates, astr)
	if err == nil {
		t.Errorf("ResolveIncludes with loop: the code did not error")
	}

}
