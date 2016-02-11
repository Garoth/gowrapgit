package gowrapgit

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
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
	if err = Clone(source, path, bare); err != nil {
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

func TestFindGits(t *testing.T) {
	t.Log("Creating fresh git clone...")

	repo1 := setupTestClone(false, t)
	location := filepath.Dir(repo1)
	repo2 := filepath.Join(location, "a", "b", "c")
	repo3 := filepath.Join(location, "flat")
	repo4 := filepath.Join(location, "xxx", "yyy")
	defer cleanupTestClone(repo1, t)

	t.Log(" - Test repo 1 cloned to", prettyPath(repo1))

	if err := Clone(repo1, repo2, true); err != nil {
		t.Fatal("Failed to clone repo:", err)
	}
	defer cleanupTestClone(repo2, t)

	t.Log(" - Test repo 2 cloned")

	if err := Clone(repo1, repo3, false); err != nil {
		t.Fatal("Failed to clone repo:", err)
	}
	defer cleanupTestClone(repo3, t)

	t.Log(" - Test repo 3 cloned")

	if err := Clone(repo2, repo4, false); err != nil {
		t.Fatal("Failed to clone repo:", err)
	}
	defer cleanupTestClone(repo4, t)

	t.Log(" - Test repo 4 cloned")

	results := FindGits(location)
	expected := []string{repo2, repo3, repo1, repo4}

	if reflect.DeepEqual(results, expected) == false {
		t.Fatal("FindGits failed. results =", results, "expected =", expected)
	}

	t.Log(" - FindGits found repos successfully! count =", len(results))
}

func TestBranch(t *testing.T) {
	t.Log("Cloning a git repo...")

	path := setupTestClone(false, t)
	defer cleanupTestClone(path, t)

	t.Log(" - Test repo cloned to", prettyPath(path))

	branch, err := CurrentBranch(path)

	if err != nil {
		t.Fatal("Couldn't check current branch:", err)
	}

	if branch != "master" {
		t.Fatal("Incorrect branch found. Expect master, found:", branch)
	}

	t.Log(" - Success checking branch:", branch)

	if err = Checkout(path, "HEAD~1"); err != nil {
		t.Fatal("Failed to checkout previous commit:", err)
	}

	t.Log(" - Success checking out HEAD~1 in repo")

	branch, err = CurrentBranch(path)

	if err != nil {
		t.Fatal("Couldn't check current branch:", err)
	}

	if branch != "HEAD" {
		t.Fatal("Incorrect ref name found. Expect HEAD, found:", branch)
	}

	t.Log(" - Success checking detached head branch:", branch)
}
