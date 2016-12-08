package run

import "sort"

// orderedSetString is an ordered set of string
type orderedSetString struct {
	// array: from index to string
	items []string
	// map: from string to index
	item2index map[string]int
}

// NewOrderedSetString return a new empty orderedSetString
func NewOrderedSetString() *orderedSetString {
	set := orderedSetString{
		item2index: make(map[string]int),
	}
	return &set
}

// Add method adds the given string to the set, if not already present.
// Returns the item's position.
func (set *orderedSetString) Add(s string) int {
	idx, ok := set.item2index[s]
	if !ok {
		idx = len(set.items)
		set.item2index[s] = idx
		set.items = append(set.items, s)
	}
	return idx
}

// AddSlice adds all element of the slice.
func (set *orderedSetString) AddSlice(as []string) {
	for _, s := range as {
		set.Add(s)
	}
}

// ToSlice returns the members of the set as a slice.
func (set *orderedSetString) ToSlice() []string {
	return set.items
}

// Contains returns whether the given item is in the set.
func (set *orderedSetString) Contains(s string) bool {
	_, ok := set.item2index[s]
	return ok
}

// Index returns the position of the given element in the set.
func (set *orderedSetString) Index(s string) (int, bool) {
	j, ok := set.item2index[s]
	return j, ok
}

// Value returns the element of the given position in the set.
func (set *orderedSetString) Value(idx int) string {
	return set.items[idx]
}

// Len returns the number of elements of the set.
func (set *orderedSetString) Len() int {
	return len(set.items)
}

// Sort order the array of string.
func (set *orderedSetString) Sort() {
	sort.Strings(set.items)
	for idx, item := range set.items {
		set.item2index[item] = idx
	}
}
