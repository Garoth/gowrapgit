package gowrapgit

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

// Returns a command object for the given shell command.
func command(name string, args ...string) *exec.Cmd {
	cmd := exec.Command(name, args...)
	// cmd.Dir = manager.root
	return cmd
}

// Checks if a file or directory exists at the given path.
func exists(path string) error {
	_, err := os.Stat(path)
	return err
}

// Checks whether this package has everything it needs to be able to function.
func sanityCheck() error {
	_, err := exec.LookPath("git")
	if err != nil {
		log.Println("Sanity Check Failure: Couldn't find 'git'.")
	}
	return err
}

// IsRepo returns whether the folder specified by the path is a git repo.
func IsRepo(path string) error {
	cmd := command("git", "rev-parse", "--git-dir")
	cmd.Dir = path
	return cmd.Run()
}

// IsBareRepo checks whether the folder specified is a git repo and bare.
func IsBareRepo(path string) error {
	var err error

	if err = sanityCheck(); err != nil {
		return err
	}

	if err = IsRepo(path); err != nil {
		return err
	}

	cmd := command("git", "rev-parse", "--is-bare-repository")
	cmd.Dir = path
	var out []byte
	out, err = cmd.Output()
	if err != nil {
		return err
	}

	if bytes.Contains(out, []byte{'t', 'r', 'u', 'e'}) == false {
		return errors.New("Not a bare repository")
	}

	return nil
}

// Clone clones a git repo. It takes a full path to clone to, a full path
// (or URL) to clone from, and whether it should be a bare git repo.
func Clone(cloneFromPath, cloneToPath string, bare bool) error {
	if err := sanityCheck(); err != nil {
		return err
	}

	// Don't check out stuff that already exists
	if exists(cloneToPath) == nil {
		return nil
	}

	if bare {
		return command("git", "clone", "--bare", cloneFromPath, cloneToPath).Run()
	}

	return command("git", "clone", cloneFromPath, cloneToPath).Run()
}

// FindGits returns a set of paths for all found git repositories under the
// given starting filepath. This will recursively walk down to find them.
func FindGits(start string) []string {
	var paths []string

	walker := func(path string, info os.FileInfo, err error) error {
		if IsRepo(path) == nil {
			paths = append(paths, path)
			return filepath.SkipDir
		}

		return nil
	}

	filepath.Walk(start, walker)

	return paths
}

// CurrentBranch reports the current active branch for the repo at the given
// path. As normal, returns "HEAD" when we're not on an otherwise named head.
func CurrentBranch(path string) (string, error) {
	if err := sanityCheck(); err != nil {
		return "", err
	}

	cmd := command("git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Dir = path
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(out)), nil
}

// MakeBranch creates a new local branch with the given `newBranchName`. If
// `sourceBranchName` is not "", then it is used as the basis for the new
// branch.
func MakeBranch(path, newBranchName, sourceBranchName string) error {
	if err := sanityCheck(); err != nil {
		return err
	}

	var cmd *exec.Cmd
	if sourceBranchName == "" {
		cmd = command("git", "branch", newBranchName)
	} else {
		cmd = command("git", "branch", newBranchName, sourceBranchName)
	}
	cmd.Dir = path
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

// ListBranches returns a list of branch strings for the git repo at
// the given path. The `local` parameter switches between local and
// remote branches.
func ListBranches(path string, local bool) ([]string, error) {
	if err := sanityCheck(); err != nil {
		return []string{}, err
	}

	var gitRefsPath string
	if local {
		gitRefsPath = filepath.Join("refs", "heads")
	} else {
		gitRefsPath = filepath.Join("refs", "remotes")
	}

	cmd := command("git", "for-each-ref", "--format='%(refname)'", gitRefsPath)
	cmd.Dir = path
	out, err := cmd.Output()
	if err != nil {
		return []string{}, err
	}

	lineBytes := bytes.Split(out, []byte{'\n'})
	lineStrings := make([]string, len(lineBytes))
	for i, byteLine := range lineBytes {
		lineStrings[i] = strings.Trim(string(bytes.TrimSpace(byteLine)), "'")
	}

	// The last "line" can be just an empty string
	if lineStrings[len(lineStrings)-1] == "" {
		lineStrings = lineStrings[:len(lineStrings)-1]
	}

	return lineStrings, nil
}

// Checkout runs the git checkout command. This can be used to switch branches
// or to check out disconnected heads.
func Checkout(path, hashish string) error {
	if err := sanityCheck(); err != nil {
		return err
	}

	cmd := command("git", "checkout", hashish)
	cmd.Dir = path
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

// Commit struct that holds the sort of data that you'd expect from git log.
type Commit struct {
	Hash, Author, AuthorEmail, ParentHash, Subject, Body string
	Timestamp                                            int
}

func (commit *Commit) String() string {
	return fmt.Sprintf("\n+ Commit: %s\n| Author: %s <%s>\n| Parent: %s\n"+
		"| Timestamp: %d\n| Subject: %s\n| Body: %s", commit.Hash, commit.Author,
		commit.AuthorEmail, commit.ParentHash, commit.Timestamp,
		commit.Subject, commit.Body)
}

// NewCommit returns a commit object for the given repo path and hashish.
func NewCommit(path, hashish string) (*Commit, error) {
	logFormat := "%H%n%an%n%ae%n%ct%n%P%n%s%n%b"
	commit := &Commit{}

	cmd := command("git", "log", "-1", "--pretty="+logFormat, hashish)
	cmd.Dir = path
	output, err := cmd.Output()
	if err != nil {
		return &Commit{}, err
	}

	lineBytes := bytes.Split(output, []byte{'\n'})
	commit.Hash = string(bytes.TrimSpace(lineBytes[0]))
	commit.Author = string(bytes.TrimSpace(lineBytes[1]))
	commit.AuthorEmail = string(bytes.TrimSpace(lineBytes[2]))
	commit.Timestamp, err = strconv.Atoi(string(bytes.TrimSpace(lineBytes[3])))
	if err != nil {
		return &Commit{}, err
	}
	commit.ParentHash = string(bytes.TrimSpace(lineBytes[4]))
	commit.Subject = string(bytes.TrimSpace(lineBytes[5]))
	commit.Body = string(bytes.TrimSpace(bytes.Join(lineBytes[6:], []byte{'\n'})))

	return commit, nil
}

// Log returns an array of Commit objects, representing the history as like
// git log. Newest commit is at index zero, oldest at the end. Leave hashish
// as a blank string to get the default "git log" result.
func Log(path, hashish string) ([]*Commit, error) {
	logFormat := "%H" // Just the hashes
	var cmd *exec.Cmd

	if hashish == "" {
		cmd = command("git", "log", "--pretty="+logFormat)
	} else {
		cmd = command("git", "log", "--pretty="+logFormat, hashish)
	}
	cmd.Dir = path
	output, err := cmd.Output()
	if err != nil {
		return []*Commit{}, err
	}

	lineBytes := bytes.Split(output, []byte{'\n'})
	// The last split is just an empty string, right?
	lineBytes = lineBytes[0 : len(lineBytes)-1]
	commits := make([]*Commit, len(lineBytes))

	for x := 0; x < len(lineBytes); x++ {
		commit, commitErr := NewCommit(path, string(lineBytes[x]))
		if commitErr != nil {
			return []*Commit{}, commitErr
		}
		commits[x] = commit
	}

	return commits, nil
}
