# Go Wrap Git

![Go Wrap Git Logo](https://garoth.github.com/img/gowrapgit-logo-small.png)

This project provides a Go API to common Git functions. It is implemented
as a simple `git` command shell wrapper, so it doesn't have any complex
non-Go dependencies. Want to know if it's any good?

- Documentation: https://godoc.org/github.com/Garoth/gowrapgit
- Code Quality Report: https://goreportcard.com/report/github.com/Garoth/gowrapgit

The aim of this project is to provide the most-used features of common commands.
This project is currently geared more towards querying git repositories rather
than affecting and synchronizing them.

If you're looking for a much more complete and down-to-the-metal implementation,
I recommend looking at https://github.com/libgit2/git2go. That project wraps
around libgit2 directly.

## Project & Testing Status

Generally pretty beta. I'm mostly just adding functions as I need them
for other projects.  However, the project is outstandingly simple, so maybe
it makes sense for you to add whatever you need.

Although this project lacks completeness, the features that it does have are
heavily unit tested. This is uniformly true for all currently implemented
functions. Currently, there is about 50% more unit test code than actual code.

## Feature Plan / Progress

Documentation: https://godoc.org/github.com/Garoth/gowrapgit

- [X] Check If Is Repo
    - [X] Check If Is Bare Repo
- [X] Log
    - [X] General Commit Struct
    - [X] Log Produces Array of Commit Structs
- [X] Clone
    - [X] --bare Option
- [X] Checkout
- [X] Branch
    - [X] Get Current Branch
    - [X] Get All Local Branches
    - [X] Get All Remote Branches
    - [X] Make Local Branch
- [X] Find All Git Repos Recursively From Path
- [X] Reset
    - [X] --hard Option
- [ ] **(future)** Status
- [ ] **(future)** Worktree 
    - [ ] Spawn New Worktree From Parent
    - [ ] Find Parent Of Worktree
- [ ] **(future)** Tree Manipulation Features
    - [ ] Fetch
    - [ ] Merge
    - [ ] Rebase
    - [ ] Push
    - [ ] Pull
    - [ ] Add
    - [ ] Commit
