package gowrapgit

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"testing"
)

func TestSanityCheck(t *testing.T) {
	t.Log("Running general sanity check...")

	if err := sanityCheck(); err != nil {
		t.Errorf("Sanity check failed, no git?")
	}
}

// Sets up a test repo that it clones from the internet and returns the path.
func setupTestClone(bare bool, t *testing.T) string {
	var err error
	var path string
	numCommits := 3

	path, err = ioutil.TempDir("", "gowrapgit")
	if path == "" || err != nil {
		t.Fatal("Couldn't create tmpdir:", err)
	}

	cmd := command("git", "init")
	cmd.Dir = path
	if err = cmd.Run(); err != nil {
		t.Fatal("Couldn't init git:", err)
	}

	for x := 0; x < numCommits; x++ {
		cmd := command("git", "commit", "--allow-empty",
			"--date", "1000000000",
			"--author", "Andrei Thorp <garoth@gmail.com>",
			"-m", "subject "+strconv.Itoa(x),
			"-m", "body message")
		cmd.Dir = path
		env := os.Environ()
		env = append(env, fmt.Sprintf("GIT_COMMITTER_DATE=%d", 1000000000))
		cmd.Env = env

		if err = cmd.Run(); err != nil {
			t.Fatal("Couldn't make commit git:", err)
		}
	}

	if bare {
		var path2 string
		path2, err = ioutil.TempDir("", "gowrapgit")
		path2 = path2 + "/clone"
		Clone(path, path2, true)
		defer cleanupTestClone(path, t)
		return path2
	}

	return path
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
	location, dirErr := ioutil.TempDir("", "gowrapgit")
	if location == "" || dirErr != nil {
		t.Fatal("Couldn't create tmpdir:", dirErr)
	}
	repo2 := filepath.Join(location, "a", "b", "c")
	repo3 := filepath.Join(location, "flat")
	repo4 := filepath.Join(location, "xxx", "yyy")
	defer cleanupTestClone(location, t)

	t.Log(" - Test repo 1 cloned to", prettyPath(repo1))

	if err := Clone(repo1, repo2, true); err != nil {
		t.Fatal("Failed to clone repo:", err)
	}

	t.Log(" - Test repo 2 cloned")

	if err := Clone(repo1, repo3, false); err != nil {
		t.Fatal("Failed to clone repo:", err)
	}

	t.Log(" - Test repo 3 cloned")

	if err := Clone(repo2, repo4, false); err != nil {
		t.Fatal("Failed to clone repo:", err)
	}

	t.Log(" - Test repo 4 cloned")

	results := FindGits(location)
	// We don't find the starting repo since it's not under location
	expected := []string{repo2, repo3, repo4}

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

func TestNewCommit(t *testing.T) {
	expected := &Commit{}
	expected.author = "Andrei Thorp"
	expected.authorEmail = "garoth@gmail.com"
	expected.timestamp = 1000000000
	expected.subject = "subject 2"
	expected.body = "body message"

	path := setupTestClone(false, t)
	defer cleanupTestClone(path, t)
	t.Log(" - Test repo cloned to", prettyPath(path))

	commit, err := NewCommit(path, "HEAD")
	if err != nil {
		t.Fatal("Error creating new Commit object:", err)
	}

	t.Logf(" - Created new Commit")

	// TODO: should check hashes as well, but they change...
	commit.hash = ""
	commit.parentHash = ""

	if reflect.DeepEqual(commit, expected) == false {
		t.Fatal("Commit isn't as expected!\nUnexpected commit data = ", commit,
			"\nExpected commit data =", expected)
	}

	t.Log(" - New Commit matches expectations")
}

func TestLog(t *testing.T) {
	path := setupTestClone(false, t)
	defer cleanupTestClone(path, t)
	t.Log(" - Test repo cloned to", prettyPath(path))

	log, err := Log(path, "")
	if len(log) == 0 || err != nil {
		t.Fatal("Couldn't get log. err:", err, "| log:", log)
	}

	t.Log(" - Success getting log:", log)

	// TODO: verify log results more, but it's sorta covered
}

// TODO: test log starting at various hashes
// func TestLogHash(t *testing.T) {
// }
