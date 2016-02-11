package gowrapgit

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestSanityCheck(t *testing.T) {
	t.Log("Running general sanity check...")

	if err := sanityCheck(); err != nil {
		t.Errorf("Sanity check failed, no git?")
	}
}

// Sets up a test repo that it clones from the internet and returns the path.
func setupTestClone(bare bool, t *testing.T) (path string) {
	var err error
	source := "https://github.com/Garoth/go-signalhandlers"

	path, err = ioutil.TempDir("", "gowrapgit")
	if path == "" || err != nil {
		t.Fatal("Couldn't create tmpdir:", err)
	}

	path = path + "/go-signalhandlers"
	if err = Clone(path, source, bare); err != nil {
		t.Fatal("Error cloning go-signalhandlers:", err)
	}

	return
}

// Deletes the repo given by the path.
func cleanupTestClone(path string, t *testing.T) {
	if err := os.RemoveAll(path); err != nil {
		t.Log("Couldn't remove test clone:", err)
	}
}

// Returns a pretty version of a given path.
func prettyPath(path string) string {
	last := filepath.Base(path)
	secondLast := filepath.Base(filepath.Dir(path))
	if last != "" && secondLast != "" {
		return filepath.Join("...", last, secondLast)
	}

	return "ERROR PRETTYING PATH!"
}

func TestClone(t *testing.T) {
	t.Log("Cloning a git repo...")

	path := setupTestClone(false, t)
	defer cleanupTestClone(path, t)

	t.Log(" - Test repo cloned to", prettyPath(path))

	if IsRepo(path) != nil {
		t.Fatal("Newly cloned repo doesn't exist?")
	}

	t.Log(" - IsRepo confirmed success")

	if IsBareRepo(path) == nil {
		t.Fatal("Newly cloned repo is bare?")
	}

	t.Log(" - IsBareRepo confirmed repo as NOT BARE")
}

func TestCloneBare(t *testing.T) {
	t.Log("Cloning a git repo...")

	path := setupTestClone(true, t)
	defer cleanupTestClone(path, t)

	t.Log(" - Test repo cloned to", prettyPath(path))

	if IsRepo(path) != nil {
		t.Fatal("Newly cloned repo doesn't exist?")
	}

	t.Log(" - IsRepo confirmed success")

	if err := IsBareRepo(path); err != nil {
		t.Fatal("Newly cloned repo isn't bare?", err)
	}

	t.Log(" - IsBareRepo confirmed repo as BARE")

}
