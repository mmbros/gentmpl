package run

import (
	"bytes"
	"fmt"
	"strconv"
)

// usize returns the number of bits of the smallest unsigned integer
// type that will hold n. Used to create the smallest possible slice of
// integers to use as indexes into the concatenated strings.
func usize(n int) int {
	switch {
	case n < 1<<8:
		return 8
	case n < 1<<16:
		return 16
	default:
		// 2^32 is enough constants for anyone.
		return 32
	}
}

// return "uint8" | "uint16" | "uint32" based on `n`
func uint(n int) string {
	return "uint" + strconv.Itoa(usize(n))
}

// astr2str returns a string representation of the items.
// Example: astr2str([]string{"a", "b", "c"}) -> "\"a\", \"b\", \"c\""
func astr2str(items []string) string {
	b := new(bytes.Buffer)

	for j, s := range items {
		if j > 0 {
			fmt.Fprint(b, ", ")
		}
		fmt.Fprintf(b, "%q", s)
	}
	return b.String()
}

// aint2str returns a string representation of the items.
// Example: aint2str([]int{1, 2, 3}) -> "1, 2, 3"
func aint2str(items []int) string {
	b := new(bytes.Buffer)

	for j, item := range items {
		if j > 0 {
			fmt.Fprint(b, ", ")
		}
		fmt.Fprintf(b, "%d", item)
	}
	return b.String()
}
