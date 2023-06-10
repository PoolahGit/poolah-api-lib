# Overview
This is the repo that will handle creating our most re-used boilerplate functions we need for API development

1.) AWS config intialization
2.) DB initialization
3.) Middleware Logic for authentication

## Prereqs
1.) Have go installed!

```bash
brew update && brew install golang
```
* in your zprofile:

``` bash
export GOPATH=$HOME/go
export GOROOT="$(brew --prefix golang)/libexec"
export PATH="$PATH:${GOPATH}/bin:${GOROOT}/bin"
```

* rerun your profile:
```bash
source $HOME/.bashrc
```
