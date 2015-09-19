package version

import (
	"os"
	"testing"

	"github.com/fern4lvarez/piladb/pkg/version/_vendor/src/github.com/mitchellh/go-homedir"
)

func TestCommitHash(t *testing.T) {
	if len(CommitHash()) != 41 {
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
	if h := CommitHash(); h != "unknown" {
		t.Errorf("CommitHash is %s, expected %s", h, "unknown")
	}

}
