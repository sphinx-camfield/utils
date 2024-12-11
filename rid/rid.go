package rid

import (
	"encoding/base64"
	"fmt"
	"github.com/google/uuid"
	"strconv"
	"strings"
	"time"
)

// Rid represents a resource identifier, which is a unique one among all resources.
// The Rid will have less than 128 characters, URL-friendly, unique, and sortable.
//
// It is composed of three parts: resource name, index, and unique identifier.
// - The resource name is a string that represents the resource type.
// - The index is a hex string that represents the time when the resource is created. It will increase over time.
// - The unique identifier is a hex string that represents a random number.
type Rid struct {
	res string
	idx string
	uni string
}

// New creates a new Rid with the given resource name.
func New(res string) Rid {
	// idx: time in hex
	idx := strconv.FormatInt(time.Now().UnixNano(), 16)
	// uni: uuid in base64
	uni, _ := strings.CutSuffix(base64.URLEncoding.EncodeToString([]byte(uuid.NewString())), "=")
	return Rid{
		res: res,
		idx: idx,
		uni: uni,
	}
}

// String returns the string representation of the Rid.
func (rid Rid) String() string {
	return rid.res + "." + rid.idx + "." + rid.uni
}

// Parse parses the given string and returns the Rid.
func Parse(s string, res string) (Rid, error) {
	parts := strings.Split(s, ".")
	rid := Rid{}

	if len(parts) != 3 {
		return rid, fmt.Errorf("invalid rid")
	}

	if parts[0] != res {
		return rid, fmt.Errorf("invalid resource name")
	}

	idxTs, err := strconv.ParseInt(parts[1], 16, 64)
	if err != nil {
		return rid, fmt.Errorf("invalid index")
	}

	if time.Now().UnixNano()-idxTs < 0 {
		// The index is in the future.
		return rid, fmt.Errorf("invalid index")
	}

	decodeString, err := base64.URLEncoding.DecodeString(parts[2])
	if err != nil {
		return rid, fmt.Errorf("invalid unique identifier")
	}

	if _, err = uuid.Parse(string(decodeString)); err != nil {
		return rid, fmt.Errorf("invalid unique identifier")
	}

	rid.res = parts[0]
	rid.idx = parts[1]
	rid.uni = parts[2]

	return rid, nil
}
