package run

import "testing"

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
