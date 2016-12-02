package run

import "testing"

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
	res, _ := resolveIncludes(templates, astr)
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

	_, err := resolveIncludes(templates, astr)
	if err == nil {
		t.Errorf("resolveIncludes with loop: the code did not error")
	}

}

/*
func oldTestResolveIncludesPanic(t *testing.T) {
	// http://stackoverflow.com/questions/31595791/how-to-test-panics
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("resolveIncludes with loop: the code did not panic")
		}
	}()
	var templates = map[string][]string{
		"incA": {"A1", "incB"},
		"incB": {"incA", "B1"},
	}
	var astr = []string{"incA"}

	_, _ = resolveIncludes(templates, astr)

}
*/

func TestUsize(t *testing.T) {
	var cases = []struct {
		input    int
		expected int
	}{
		{0, 8},
		{255, 8},
		{256, 16},
		{65535, 16},
		{65536, 32},
		{-1, 8},
		{-256, 8},
		{-65536, 8},
	}

	for _, c := range cases {
		actual := usize(c.input)
		if actual != c.expected {
			t.Errorf("usize(%d): expected %d, actual %d", c.input, c.expected, actual)
		}
	}

}
