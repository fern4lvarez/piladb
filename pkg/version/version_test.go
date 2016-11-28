package version

import (
	"os"
	"testing"

	"github.com/mitchellh/go-homedir"
)

func TestCommitHash(t *testing.T) {
	if len(CommitHash()) != 40 {
		t.Errorf("commit hash version unexpected")
	}
}

func TestCommitHash_Unknown(t *testing.T) {
	// normally $HOME is not version controlled
	home, err := homedir.Dir()
	if err != nil {
		t.Fatal("could not get home dir")
	}

	os.Chdir(home)
	if h := CommitHash(); h != "master" {
		t.Errorf("CommitHash is %s, expected %s", h, "unknown")
	}

}
