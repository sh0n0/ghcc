# ghcc
ghcc combines ghq (repository management tool) with scc (code counter).
## Install

```
$ go get github.com/kwQt/ghcc
```

## Dependencies
- [x-motemen/ghq](https://github.com/x-motemen/ghq)
- [boyter/scc](https://github.com/boyter/scc)

## Commands

### get
get source code from repository by using ghq get, and then show scc results.


```
$ ghcc get https://github.com/<user>/<repository>.git
```
You can remove source code automatically if you enter "y" at the last prompt.

### ls
display history of scc results as list.

```
$ ghcc ls
```

### clear
clear all history
```
$ ghcc clear
```
