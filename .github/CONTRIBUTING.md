# Contributing to go-testparrot

**First:** if you're unsure or afraid of _anything_, just ask or submit the
issue or pull request anyway. You won't be yelled at for giving your best
effort. The worst that can happen is that you'll be politely asked to change
something. Any sort of contribution is appreciated, and don't want a wall of
rules to get in the way of that.

However, for those individuals who want a bit more guidance on the best way to
contribute to the project, read on. This document will cover what we're looking
for. By addressing all the points we're looking for, it raises the chances we
can quickly merge or address your contributions.

## Issues and pull requests

### Reporting an Issue

- Make sure you test against the latest released version. It is
  possible the bug you're experiencing has already been fixed.

- Provide a reproducible test case. If a contributor can't reproduce an issue,
  then it dramatically lowers the chances it'll get fixed. And in some cases,
  the issue will eventually be closed.

- Respond promptly to any questions made by the developer team to your issue.
  Stale issues will be closed.

### Opening an Pull Request

Thank you for contributing! When you are ready to open a pull-request, you will
need to
[fork go-testparrot](https://github.com/xtruder/go-testparrot#fork-destination-box,
push yourchanges to your fork, and then open a pull-request.

For example, my github username is `offlinehacker`, so I would do the following:

```
git checkout -b f-my-feature
# Develop a patch.
git push https://github.com/offlinehacker/go-testparrot f-my-feature
```

From there, open your fork in your browser to open a new pull-request.

**Note:** See '[Working with
forks](https://help.github.com/articles/working-with-forks/)' for a better way
to use `git push ...`.

## Development

This project provides multiple ways to develop project, and not forcing
you in a speciffic development flow. In general having latest go is enough,
but using [vscode remote containers](https://code.visualstudio.com/docs/remote/containers) or [nix-shell](https://nixos.org/), which provides deterministic
development environments is advised.

### Using visual studio code development containers

If you are using [visual studio code](https://code.visualstudio.com/) you can open
project in [vscode development container](https://code.visualstudio.com/docs/remote/containers). This allows you to develop your project in container in
deterministic way.

First you will need to install [vscode remote extensionpack](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.vscode-remote-extensionpack).

After that you can open project in container using these [instructions](https://code.visualstudio.com/docs/remote/containers#_quick-start-open-an-existing-folder-in-a-container).

### Using nix-shell for development

If you want to have deterministic development environment, but don't like to use
vscode remote containers, you can also use [nix-shell](https://nixos.org/) for
development. Make sure to [install nix](https://nixos.org/download.html) first.

After you have nix installed you can simply use `nix-shell` command to enter
development environment.

### Using go for develoment

If you have never worked with Go before, you will have to install its
runtime in order to develop go-testparrot.

1. This project always releases from the latest stable version of golang.
   [Install go](https://golang.org/doc/install#install) 

2. Make sure go modules are enabled by setting `GO111MODULE=on` environment
   variable.

### Running tests

```bash
go test ./
```
