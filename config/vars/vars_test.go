package vars

import "testing"

func TestEnv(t *testing.T) {
	name := "FOO"
	expectedEnv := "PILADB_FOO"

	if e := Env(name); e != expectedEnv {
		t.Errorf("Env is %s, expected %s", e, expectedEnv)
	}
}
