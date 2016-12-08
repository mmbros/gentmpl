package run

import (
	"fmt"
	"testing"
)

func TestOSS(t *testing.T) {

	oss := NewOrderedSetString()

	oss.AddSlice([]string{"item1", "item2", "item3"})
	if expect := 3; oss.Len() != expect {
		t.Errorf("AddSlice: Len: expecting %d, found %d", expect, oss.Len)
	}

	oss.AddSlice([]string{"item2", "item4", "item1"})
	if expect := 4; oss.Len() != expect {
		t.Errorf("AddSlice: Len: expecting %d, found %d", expect, oss.Len)
	}

	for expect := 0; expect < 4; expect++ {
		item := fmt.Sprintf("item%d", expect+1)
		pos, ok := oss.Index(item)
		if !ok {
			t.Errorf("Index(%q): expecting %d, but not found", item, expect)
		} else if pos != expect {
			t.Errorf("Index(%q): expecting %d, found %d", item, expect, pos)
		}
	}

}
