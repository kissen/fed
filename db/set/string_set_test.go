package set

import (
	"strings"
	"testing"
)

func TestStringSet(t *testing.T) {
	set := NewStringSet()

	words := strings.Split(
		"thy brother hecuba from dymas sprung a valiant warrior haughty bold and young",
		" ",
	)

	// add

	for _, word := range words {
		if added := set.Put(word); !added {
			t.Errorf("failed to add word=%v", word)
		}
	}

	// len

	if actual := set.Len(); actual != len(words) {
		t.Errorf("bad return of Len() expected=%v actual=%v", len(words), actual)
	}

	// get

	for _, word := range words {
		if !set.Contains(word) {
			t.Errorf("word=%v missing even though it should have been added", word)
		}
	}

	// remove existing

	if removed := set.Remove("valiant"); !removed {
		t.Errorf("refusing to remove added entry")
	}

	// try removing missing

	if removed := set.Remove("hector"); removed {
		t.Errorf("confirming removal of item that was not added")
	}
}
