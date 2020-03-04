package ap

import "testing"

func TestIsEmpty_Empty(t *testing.T) {
	if !isEmpty("") {
		t.Error("empty string is judged not empty")
	}
}

func TestIsEmpty_Whitespace(t *testing.T) {
	if !isEmpty(" ") {
		t.Error("whitespace string is judged not empty")
	}

	if !isEmpty("\t") {
		t.Error("whitespace string is judged not empty")
	}

	if !isEmpty("\n") {
		t.Error("whitespace string is judged not empty")
	}
}

func TestIsEmpty_NotEmpty(t *testing.T) {
	if isEmpty("fckafd") {
		t.Errorf("not empty string is judged empty")
	}
}
