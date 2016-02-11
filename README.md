# Go Wrap Git

This project provides a Go API to common Git functions. It is implemented
as a simple `git` command shell wrapper, so it doesn't have any complex
non-Go dependencies. Want to know if it's any good?

- Documentation: https://godoc.org/github.com/Garoth/gowrapgit
- Code Quality Report: https://goreportcard.com/report/github.com/Garoth/gowrapgit

If you're looking for a much more complete and down-to-the-metal implementation,
I recommend looking at https://github.com/libgit2/git2go. That project wraps
around libgit2 directly.

## Project Status

Generally pretty beta. I'm mostly just adding functions as I need them
for other projects. You can expect it to have just the bare minimum.
However, the project is braindead simple, so maybe it makes sense for you
to add whatever you need.

### Testing Status

Although this project lacks completeness, the features that it does have are
heavily unit tested. This is uniformly true for all currently implemented
functions.

### Feature Plan / Progress

Documentation: https://godoc.org/github.com/Garoth/gowrapgit

- [X] Check If Is Repo
    - [X] Check If Is Bare Repo
- [ ] Log
    - [X] General Commit Struct That Can Hold Relevant Data
    - [ ] Hook Up Log To Produce Array of Commit Structs
    - [ ] Unit Test Using Locally Made Repo
- [X] Clone
    - [X] Clone Bare
- [X] Checkout
- [ ] Branch
    - [X] Get Current Branch
    - [ ] Get All Local Branches
    - [ ] Get All Remote Branches
- [X] Find All Git Repos Recursively From Path
- [ ] Reset (hard and soft)
- [ ] Worktree 
    - [ ] Spawn New Worktree From Parent
    - [ ] Find Parent Of Worktree
- [ ] Tree Manipulation Features **(future)**
    - [ ] Fetch
    - [ ] Merge
    - [ ] Rebase
    - [ ] Push
    - [ ] Pull
    - [ ] Add
    - [ ] Commit
