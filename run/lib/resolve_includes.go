// Package lib implements some library functions used by gentmpl utility.
package lib

import "fmt"

// ResolveIncludes returns a modified version of the input mapping.
// It filter the elements taking only the keys found in the names array.
// Moreover every item that corrispond to a mapping key name is (recursivery)
// expanded with the mapping items.
// It returs an error in case of cyclic includes.
func ResolveIncludes(mapping map[string][]string, names []string) (map[string][]string, error) {
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
