package rids

import "testing"

func TestRids_MakeValid(t *testing.T) {
	id, err := Make("test")

	if err != nil {
		t.Error(err)
	}

	if id == "" {
		t.Error("id is empty")
	}

	valid, err := Valid(id)
	if err != nil {
		t.Error(err)
	}

	if !valid {
		t.Error("id should be valid")
	}
}
