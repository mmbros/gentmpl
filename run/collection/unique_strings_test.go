package collection

import (
	"fmt"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestUniqueStrings(t *testing.T) {

	checkLen := func(expect, actual int) {
		if actual != expect {
			t.Errorf("Len: expecting %d, found %d", expect, actual)
		}
	}
	checkIdx := func(item string, expect, actual int) {
		if actual != expect {
			t.Errorf("Index(%q): expecting %d, found %d", item, expect, actual)
		}
	}

	var idx int
	col := NewUniqueStrings()

	// initial len
	checkLen(0, col.Len())

	// add item 1
	idx = col.Add("item1")
	checkIdx("item1", 0, idx)
	checkLen(1, col.Len())

	// add item 2
	idx = col.Add("item2")
	checkIdx("item2", 1, idx)
	checkLen(2, col.Len())

	// add item 1 (already existing)
	idx = col.Add("item1")
	checkIdx("item1", 0, idx)
	checkLen(2, col.Len())

	// add slice of items
	col.AddSlice([]string{"item1", "item2", "item3", "item4"})
	checkLen(4, col.Len())

	// check index
	for expect := 0; expect < 4; expect++ {
		item := fmt.Sprintf("item%d", expect+1)
		pos, ok := col.Index(item)
		if !ok {
			t.Errorf("Index(%q): expecting %d, but not found", item, expect)
		} else {
			checkIdx(item, expect, pos)
		}
	}

	// check value
	for idx := 0; idx < 4; idx++ {
		expect := fmt.Sprintf("item%d", idx+1)
		actual := col.Value(idx)
		if actual != expect {
			t.Errorf("Value(%d): expecting %q, found %q", idx, expect, actual)
		}
	}

}

func TestUniqueStrings_ErriIndexOutOfRange(t *testing.T) {
	defer func() {
		msg := "runtime error: index out of range"
		showerr := func(txt string) {
			t.Errorf("Expecting %q, got: %q", msg, txt)
		}
		r := recover()
		if r == nil {
			showerr("no error")
		} else {
			e, ok := r.(error)
			// TODO: find a better to check "index out of range" error
			if ok {
				if strings.Index(e.Error(), msg) != 0 {
					showerr(e.Error())
				}
			} else {
				showerr(fmt.Sprintf("%v", r))
			}
		}
	}()
	col := NewUniqueStrings()
	col.AddSlice([]string{"item1", "item2", "item3", "item4"})
	// check value with index out of range
	_ = col.Value(4)
}
func TestUniqueStrings_Sort(t *testing.T) {
	myslice := []string{"item3", "item2", "item4", "item1"}
	col := NewUniqueStrings()
	col.AddSlice(myslice)

	// check initial position
	itemPos := map[string]int{}
	for i, k := range myslice {
		itemPos[k] = i
	}
	for actual, k := range myslice {
		expect := itemPos[k]
		if actual != expect {
			t.Errorf("Index(%q): expecting %d, found %d", k, expect, actual)
		}
	}

	// sort
	col.Sort()

	// check sorted positions
	for idx := 0; idx < 4; idx++ {
		expect := fmt.Sprintf("item%d", idx+1)
		actual := col.Value(idx)
		if actual != expect {
			t.Errorf("Value(%d): expecting %q, found %q", idx, expect, actual)
		}
	}
}

func TestUniqueStrings_Contains(t *testing.T) {
	myslice := []string{"item3", "item2", "item4", "item1", "item3", "item3"}
	col := NewUniqueStrings()
	col.AddSlice(myslice)

	tests := []struct {
		item string
		want bool
	}{
		{item: "item1", want: true},
		{item: "item2", want: true},
		{item: "item3", want: true},
		{item: "item4", want: true},
		{item: "item5", want: false},
		{item: "item6", want: false},
	}

	for _, test := range tests {
		got := col.Contains(test.item)
		if got != test.want {
			t.Errorf("Contains(%q) = %v, want %v", test.item, got, test.want)
		}
	}
}

func TestUniqueStrings_ToSlice(t *testing.T) {
	myslice := []string{"item3", "item2", "item3", "item4", "item1", "item3"}
	col := NewUniqueStrings()
	col.AddSlice(myslice)

	want := []string{"item3", "item2", "item4", "item1"}
	got := col.ToSlice()

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("ToSlice() mismatch (-want +got):\n%s", diff)
	}

}
