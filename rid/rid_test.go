package rid

import (
	"fmt"
	"strings"
	"testing"
)

func TestRid_String(t *testing.T) {
	r := New("test")

	fmt.Println("rid: ", r.String(), len(r.String()))

	if strings.HasPrefix(r.String(), "test.") == false {
		t.Error("strings.HasPrefix(r.String(), \"test.\") == false")
	}

	rid, err := Parse(r.String(), "test")

	if err != nil {
		t.Error(err)
	}

	if r.String() != rid.String() {
		t.Error("rid.String() != rid.String()")
	}
}
