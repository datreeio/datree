# Developer Guide

This guide explains how to set up your environment for developing on Datree.  
This guide was written for macOS and Linux machines.

## Prerequisites

- Go version 1.18
- Git

## Building Datree

We use [Make](https://www.gnu.org/software/make/) to build our programs. The simplest way to get started is:

```
$ make build
```

This will build the executable file and place it in the project root.

One way to run datree locally is to use the newly created executable file:

```
$ ./datree test ./internal/fixtures/kube/k8s-demo.yaml
```

## Running tests

To run all the tests:

```
$ make test
```

## Contribution Guidelines

Make sure you have read and understood the
main [CONTRIBUTING](https://github.com/datreeio/datree/blob/main/CONTRIBUTING.md) guide:

### Structure of the Code

#### cobra

We use [cobra](https://github.com/spf13/cobra) as our Command Line Interface framework.  
The available commands can be found under the `cmd` directory; each folder represents a datree command. To add a
command, add a folder with the cobra command, and use it in the `cmd/root.go` file.

#### api endpoints

Datree requires an internet connection to connect to our backend API.  
While developing locally, API requests will reach our staging environment and be visible on
the [Staging Dashboard](https://app.staging.datree.com).  
All available API requests can be found under `pkg/cliClient`

#### manual testing

It's best to use fixtures for manual testing, which are found under `internal/fixtures`

#### test coverage

To add a test for a given file, add a file with a `_test` suffix.  
For example: for the file `./reader.go` add a test file `./reader_test.go`

- For bug fixes: add a test that covers the bug fixed
- For features: add tests for the feature

### Git Conventions

The main branch is the home of the current development candidate.  
We accept changes to the code via GitHub Pull Requests (PRs). One workflow for doing this is as follows:

1. Fork the repository and clone it locally.
2. Create a new working branch `git checkout -b "ISSUE#195_some_short_description"`
3. When you are ready for us to review your changes, push your branch to GitHub, and then open a new pull request to
   the `main` branch.

For Git commit messages, please follow
our [Commit Message Format](https://github.com/datreeio/datree/blob/main/CONTRIBUTING.md#-commit-message-format).  
For example: `git commit -m "feat: add windows support"`

### Go Conventions

We follow the [standard go formatting](https://golang.org/doc/effective_go#formatting) - simply use your IDE's
auto-formatter to make your code beautiful.
