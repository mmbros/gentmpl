package run

import "sort"

// OrderedSetString is an ordered set of string
type OrderedSetString struct {
	items      []string
	item2index map[string]int
}

// NewOrderedSetString return a new empty OrderedSetString
func NewOrderedSetString() *OrderedSetString {
	set := OrderedSetString{
		item2index: make(map[string]int),
	}
	return &set
}

// Add adds the element. Returns the index of the element
func (set *OrderedSetString) Add(s string) int {
	idx, ok := set.item2index[s]
	if !ok {
		idx = len(set.items)
		set.item2index[s] = idx
		set.items = append(set.items, s)
	}
	return idx
}

// AddSlice adds all element of the slice.
func (set *OrderedSetString) AddSlice(as []string) {
	for _, s := range as {
		set.Add(s)
	}
}

// ToSlice returns the members of the set as a slice.
func (set *OrderedSetString) ToSlice() []string {
	return set.items
}

// Contains returns whether the given item is in the set.
func (set *OrderedSetString) Contains(s string) bool {
	_, ok := set.item2index[s]
	return ok
}

// Index returns the position of the given element in the set
func (set *OrderedSetString) Index(s string) (int, bool) {
	j, ok := set.item2index[s]
	return j, ok
}

// Value returns the element of the given position in the set
func (set *OrderedSetString) Value(idx int) string {
	return set.items[idx]
}

// Len returns the number of elements of the set
func (set *OrderedSetString) Len() int {
	return len(set.items)
}

// Sort order the array of string
func (set *OrderedSetString) Sort() {
	sort.Strings(set.items)
	for idx, item := range set.items {
		set.item2index[item] = idx
	}
}
