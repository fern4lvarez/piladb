package version

import (
	"os"
	"testing"

	"github.com/mitchellh/go-homedir"
)

func TestVersion(t *testing.T) {
	v := "1.0"
	expectedVersion := "1.0"
	if version := Version(v); version != expectedVersion {
		t.Errorf("version is %s, expected %v", version, expectedVersion)
	}

	v = ""
	if version := Version(v); len(version) != 40 {
		t.Errorf("version is unexpected: %s", version)
	}
}

func TestCommitHash(t *testing.T) {
	if c := CommitHash(); len(c) != 40 {
		t.Errorf("commit hash version unexpected: %s", c)
	}
}

func TestCommitHash_Undefined(t *testing.T) {
	// normally $HOME is not version controlled
	home, err := homedir.Dir()
	if err != nil {
		t.Fatal("could not get home dir")
	}

	os.Chdir(home)
	if h := CommitHash(); h != "undefined" {
		t.Errorf("CommitHash is %s, expected %s", h, "undefined")
	}

}
