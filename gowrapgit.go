package gowrapgit

import (
	"bytes"
	"errors"
	"log"
	"os"
	"os/exec"
	"path/filepath"
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
// path. If there is no active branch, no repo, or otherwise -- this errors.
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
