package collection

import "sort"

// UniqueStrings is a collection of distinct strings.
type UniqueStrings struct {
	// array: from index to string
	items []string
	// map: from string to index
	item2index map[string]int
}

// NewUniqueStrings return a new empty UniqueStrings collection.
func NewUniqueStrings() *UniqueStrings {
	col := UniqueStrings{
		item2index: make(map[string]int),
	}
	return &col
}

// Add method adds the given string to the collection, if not already present.
// New strings are added at the end of the collection.
// Returns the item's position.
func (col *UniqueStrings) Add(s string) int {
	idx, ok := col.item2index[s]
	if !ok {
		idx = len(col.items)
		col.item2index[s] = idx
		col.items = append(col.items, s)
	}
	return idx
}

// AddSlice adds all element of the slice to the collection.
func (col *UniqueStrings) AddSlice(astr []string) {
	for _, s := range astr {
		col.Add(s)
	}
}

// ToSlice returns the collection items as a slice.
func (col *UniqueStrings) ToSlice() []string {
	return col.items
}

// Contains returns whether the given item is in the collection.
func (col *UniqueStrings) Contains(s string) bool {
	_, ok := col.item2index[s]
	return ok
}

// Index returns the position of the given element in the collection.
func (col *UniqueStrings) Index(s string) (int, bool) {
	j, ok := col.item2index[s]
	return j, ok
}

// Value returns the element of the given position in the collection.
// if i<0 or i>=col.Len(), panic with "index out of range" runtime error.
func (col *UniqueStrings) Value(i int) string {
	return col.items[i]
}

// Len returns the number of elements of the collection.
func (col *UniqueStrings) Len() int {
	return len(col.items)
}

// Sort order the collection of strings.
func (col *UniqueStrings) Sort() {
	sort.Strings(col.items)
	for idx, item := range col.items {
		col.item2index[item] = idx
	}
}
