# Go Wrap Git

This project provides a Go API to common Git functions. It is implemented
as a simple `git` command shell wrapper, so it doesn't have any complex
dependencies.

If you're looking for a much more complete and down-to-the-metal implementation,
I recommend looking at https://github.com/libgit2/git2go. This project wraps
around libgit2 directly.

## Project Status

Generally pretty alpha. I'm basically just adding functions as I need them
for other projects. You can expect it to have just the bare minimum for now.

### Testing Status

Although this project lacks completeness, the features that it does have are
heavily unit tested. This is uniformely true for all currently implemented
functions.

### Feature Plan / Progress

- [X] Check If Is Repo
    - [X] Check If Is Bare Repo
- [ ] Log (dumped into a big struct array)
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
