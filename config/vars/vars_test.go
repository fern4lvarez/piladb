package vars

import "testing"

func TestEnv(t *testing.T) {
	name := "FOO"
	expectedEnv := "PILADB_FOO"

	if e := Env(name); e != expectedEnv {
		t.Errorf("Env is %s, expected %s", e, expectedEnv)
	}
}

func TestDefaultInt(t *testing.T) {
	inputOutput := []struct {
		input  string
		output int
	}{
		{MaxStackSize, MaxStackSizeDefault},
		{ReadTimeout, ReadTimeoutDefault},
		{WriteTimeout, WriteTimeoutDefault},
		{ShutdownTimeout, ShutdownTimeoutDefault},
		{Port, PortDefault},
		{"foo", -1},
	}

	for _, io := range inputOutput {
		if o := DefaultInt(io.input); o != io.output {
			t.Errorf("DefaultInt is %v, expected %v", o, io.output)
		}
	}
}

func TestDefaultBool(t *testing.T) {
	inputOutput := []struct {
		input  string
		output bool
	}{
		{PushWhenFull, PushWhenFullDefault},
		{NoDonate, NoDonateDefault},
		{"foo", false},
	}

	for _, io := range inputOutput {
		if o := DefaultBool(io.input); o != io.output {
			t.Errorf("DefaultBool is %v, expected %v", o, io.output)
		}
	}
}
