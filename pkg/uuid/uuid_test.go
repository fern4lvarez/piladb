package uuid

import (
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	s := "test"
	u := New(s)
	if u.String() != "1540300e31de262ee89774c014ac163d" {
		t.Errorf("u is %v, expected %v", u, "1540300e31de262ee89774c014ac163d")
	}

	u2 := New(s)
	if u != u2 {
		t.Errorf("u and u2 differ")
	}

}

func TestUUIDString(t *testing.T) {
	u := UUID("123e4567e89b12d3a456426655440000")
	s := u.String()
	if s != "123e4567e89b12d3a456426655440000" {
		t.Errorf("u.String() is %v, expected %v", s, "123e4567e89bi12d3a456426655440000")
	}

	s = fmt.Sprintf("%v", u)
	if s != "123e4567e89b12d3a456426655440000" {
		t.Errorf("u.String() is %v, expected %v", s, "123e4567e89b12d3a456426655440000")
	}
}
