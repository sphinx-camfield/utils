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

func TestRid_Collision(t *testing.T) {

	oneMillion := 1000000

	ridxChn := make(chan string, oneMillion)

	for i := 0; i < oneMillion; i++ {
		r := New("test")
		go func() {
			ridxChn <- r.String()
		}()
	}

	ridxMap := make(map[string]bool)

	for i := 0; i < oneMillion; i++ {
		id := <-ridxChn
		if _, ok := ridxMap[id]; ok {
			t.Error("collision", id)
		}
		ridxMap[id] = true
	}

	if len(ridxMap) != oneMillion {
		t.Error("len(ridxMap) != oneMillion")
	}
}
