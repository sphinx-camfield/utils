package rids

import (
	"encoding/base64"
	"fmt"
	"github.com/google/uuid"
	"strconv"
	"strings"
	"time"
)

func Make(res string) (string, error) {
	idx := strconv.FormatInt(time.Now().UnixNano(), 16)
	uni, _ := strings.CutSuffix(base64.URLEncoding.EncodeToString([]byte(uuid.NewString())), "=")
	return res + "." + idx + "." + uni, nil
}

func Valid(s string) (bool, error) {
	parts := strings.Split(s, ".")

	if len(parts) != 3 {
		return false, fmt.Errorf("invalid rid")
	}

	idxTs, err := strconv.ParseInt(parts[1], 16, 64)
	if err != nil {
		return false, fmt.Errorf("invalid index")
	}

	if time.Now().UnixNano()-idxTs < 0 {
		// The index is in the future.
		return false, fmt.Errorf("invalid index")
	}

	decodeString, err := base64.URLEncoding.DecodeString(parts[2])
	if err != nil {
		return false, fmt.Errorf("invalid unique identifier")
	}

	if _, err = uuid.Parse(string(decodeString)); err != nil {
		return false, fmt.Errorf("invalid unique identifier")
	}

	return true, nil
}
