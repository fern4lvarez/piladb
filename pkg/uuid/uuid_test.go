package uuid

import (
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	s := "test"
	u := New(s)
	if u.String() != "e8b764da-5fe5-51ed-8af8-c5c6eca28d7a" {
		t.Errorf("u is %v, expected %v", u, "e8b764da-5fe5-51ed-8af8-c5c6eca28d7a")
	}

	u2 := New(s)
	if u != u2 {
		t.Errorf("u and u2 differ")
	}

}

func TestUUIDString(t *testing.T) {
	u := UUID("4f772915-1233-5679-845f-b4fe78c3115d")
	s := u.String()
	if s != "4f772915-1233-5679-845f-b4fe78c3115d" {
		t.Errorf("u.String() is %v, expected %v", s, "4f772915-1233-5679-845f-b4fe78c3115d")
	}

	s = fmt.Sprintf("%v", u)
	if s != "4f772915-1233-5679-845f-b4fe78c3115d" {
		t.Errorf("u.String() is %v, expected %v", s, "4f772915-1233-5679-845f-b4fe78c3115d")
	}
}
