package run

import (
	"bytes"
	"fmt"
	"strconv"
)

// resolveIncludes returns a modified version of the input mapping.
// It filter the elements taking only the keys found in the names array.
// Moreover every item that corrispond to a mapping key name is (recursivery)
// expanded with the mapping items.
// It returs an error in case of cyclic includes.
func resolveIncludes(mapping map[string][]string, names []string) (map[string][]string, error) {
	type set map[string]struct{}

	m := make(map[string][]string)

	var resolve func(string, set) error

	resolve = func(name string, visited set) error {

		if _, ok := m[name]; ok {
			// already resolved
			return nil
		}

		if _, ok := visited[name]; ok {
			return fmt.Errorf("Found invalid cycle (%s)", name)
		}

		// add name to the set of already included templates
		visited[name] = struct{}{}

		// iter over each template files
		var files []string

		for _, item := range mapping[name] {
			// check if it's an include item
			if _, ok := mapping[item]; ok {
				// it's an include
				if err := resolve(item, visited); err != nil {
					return err
				}
				files = append(files, m[item]...)
			} else {
				// append the file
				files = append(files, item)
			}

		}

		m[name] = files
		return nil
	}

	res := make(map[string][]string)

	for _, name := range names {
		if err := resolve(name, set{}); err != nil {
			return nil, err
		}
		res[name] = m[name]
	}

	return res, nil
}

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
