package stringutil

import "testing"

func TestLowerFirst(t *testing.T) {
	if r := LowerFirst("Abcde"); r != "abcde" {
		t.Errorf("should be abcde but %s", r)
	}
	if r := LowerFirst(""); r != "" {
		t.Errorf("should be empty but %s", r)
	}
}
