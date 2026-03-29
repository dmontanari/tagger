# Tagger: The Git Tag Swiss Army Knife

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

. : Pragmatic Semantic Versioning for Git-based CI/CD : .

![tagger](images/tagger.png)

tagger is a CLI tool written in Go designed to automate Semantic Versioning (SemVer) based purely on Git Tags.
Motivation

In modern CI/CD ecosystems, managing artifact versions (Docker images, Binaries, Releases) has become unnecessarily bloated. Current solutions often rely on external states (Redis, databases), duplicated version files within the repo, or—worse—extracting metadata from commit messages cluttered with emojis and redundant prefixes.

tagger is built on the premise that Git is the single source of truth. Tags are native, immutable markers; using them to manage versions is the cleanest and most deterministic way to orchestrate a GitOps pipeline.

### What Tagger Solves

- Zero External Dependencies: No need for external databases or .version files that trigger circular commits in your pipeline.

- Pragmatic Logs: Frees your commit history for its original purpose: documenting technical intent and solutions, not carrying infrastructure flags.

- SemVer Consistency: Ensures that version increments (Major, Minor, Patch) strictly follow the SemVer 2.0.0 specification.

- CI/CD Performance: As a static Go binary, tagger is ideal for lightweight runners, requiring no heavy runtimes (Python, Node, PHP).

### Engineering Principles

- Single Source of Truth: The current version is always the highest SemVer tag found in the Git history.

- Immutability: Once pushed to remote, a tag is the ultimate release authority.

- Zero Bloat: Focused on doing one thing—managing tags—with maximum efficiency.

### Usage

#### Help

```bash
$ tagger help

	. : Git tag Swiss Army Knife : .

Usage:
  tagger [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  inc         Create new tag incrementing version number.
  last        Return last tag in repository path.
  list        List all tags in repository path.
  version     Show version.

Flags:
  -h, --help      help for tagger

Use "tagger [command] --help" for more information about a command.
```

#### List all tags

```bash
$ tagger list -h
list [repository path] List all tags in repository path.

Usage:
  tagger list [repository path] [flags]

Flags:
  -h, --help  help for list

$ tagger list /path/to/repo
Tags:
2026-03-28 12:26  v0.0.1
```

#### Increment Version (Major, Minor, or Patch)

tagger identifies the latest version, applies the logical increment, and optionally pushes to remote.

Note: Incrementing a higher version level resets lower ones (e.g., a Major bump on v2.1.35 results in v3.0.0).

```bash
$ tagger inc -M . --dry-run --verbose
v0.0.1 -> v1.0.0
$ tagger inc -m . --dry-run --verbose
v0.0.1 -> v0.1.0
$ tagger inc -p . --dry-run --verbose
v0.0.1 -> v0.0.2
$ tagger inc -p . --dry-run
v0.0.2
```


## Building from sources

Building from source

1. Clone this repository: git clone https://github.com/dmontanari/tagger.git

2. Build the binary: make build

The binary will be in the root directory.


## Installation

Coming soon


## License

Distributed under MIT license. See LICENSE for more information.

© 2026 Daniel Montanari. All rights reserved.


