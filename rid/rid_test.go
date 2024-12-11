package rid

import (
	"encoding/json"
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

	rid, err := Parse(r.String())

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

func TestScan(t *testing.T) {
	rid := New("test")
	ridStr := rid.String()

	var oRid *Rid

	err := Scan(ridStr, oRid)

	if err != nil {
		t.Error(err)
	}

	if rid != oRid {
		t.Error("r.String() != rid")
	}
}

func TestMust(t *testing.T) {
	rid := New("test")
	ridStr := rid.String()

	oRid := Must(ridStr)

	if rid != oRid {
		t.Error("r.String() != rid")
	}
}

func TestRid_MarshalJSON(t *testing.T) {

	rid := New("test")
	ridStr := rid.String()

	b, err := json.Marshal(&rid)

	if err != nil {
		t.Error(err)
	}

	fmt.Println("b: ", string(b))

	if string(b) != "\""+ridStr+"\"" {
		t.Error("string(b) != `\"`+ridStr+`\"`")
	}

	var oRid *Rid

	err = json.Unmarshal(b, &oRid)

	if err != nil {
		t.Error(err)
	}

	if rid != oRid {
		t.Error("Unmarshal(rid) != rid")
	}
}
