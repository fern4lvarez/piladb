package uuid

import (
	"encoding/hex"
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	s := "test"
	u := New(s)
	if u.String() != "343345ce-8555-1655-db33-5945ace62f1c" {
		t.Errorf("u is %v, expected %v", u, "343345ce-8555-1655-db33-5945ace62f1c")
	}

	u2 := New(s)
	if u != u2 {
		t.Errorf("u and u2 differ")
	}

}

func TestUUIDString(t *testing.T) {
	u := UUID("68dc78ce-5de1-11e7-907b-a6006ad3dba0")
	s := u.String()
	if s != "68dc78ce-5de1-11e7-907b-a6006ad3dba0" {
		t.Errorf("u.String() is %v, expected %v", s, "68dc78ce-5de1-11e7-907b-a6006ad3dba0")
	}

	s = fmt.Sprintf("%v", u)
	if s != "68dc78ce-5de1-11e7-907b-a6006ad3dba0" {
		t.Errorf("u.String() is %v, expected %v", s, "68dc78ce-5de1-11e7-907b-a6006ad3dba0")
	}
}

func TestCanonical(t *testing.T) {
	b, _ := hex.DecodeString("e494f6205de211e7907ba6006ad3dba0")
	expectedUUID := "e494f620-5de2-11e7-907b-a6006ad3dba0"

	s := Canonical(b)
	if s != expectedUUID {
		t.Errorf("Canonical UUID is %s, expected %s", s, expectedUUID)
	}
}

func TestCanonical_Short(t *testing.T) {
	b := []byte("abcdefgh")
	expectedUUID := "00000000-0000-0000-6162-636465666768"

	s := Canonical(b)
	if s != expectedUUID {
		t.Errorf("Canonical UUID is %s, expected %s", s, expectedUUID)
	}
}

func TestCanonical_Long(t *testing.T) {
	b, _ := hex.DecodeString("e494f6205de211e7907ba6006ad3dba0")
	b = append(b, []byte("xxxxxxxxxxxxxxxxxx")...)
	expectedUUID := "e494f620-5de2-11e7-907b-a6006ad3dba0"

	s := Canonical(b)
	if s != expectedUUID {
		t.Errorf("Canonical UUID is %s, expected %s", s, expectedUUID)
	}
}
