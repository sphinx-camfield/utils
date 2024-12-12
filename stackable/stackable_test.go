package stackable

import (
	"net/http"
	"testing"
)

func TestStackable(t *testing.T) {

	var trail []string

	s1 := func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			trail = append(trail, "s1")
			next(w, r)
		}
	}

	s2 := func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			trail = append(trail, "s2")
			next(w, r)
		}
	}

	stacked := Stackup(s1, s2)

	stacked(NoopHandlerFunc)(nil, nil)

	if len(trail) != 2 {
		t.Errorf("invalid trail, expected 2 stackables have been called")
	}

	if trail[0] != "s1" {
		t.Errorf("invalid trail, expected s1 to be called first")
	}

	if trail[1] != "s2" {
		t.Errorf("invalid trail, expected s2 to be called second")
	}
}
