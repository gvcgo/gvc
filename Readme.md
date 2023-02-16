## What is gvc?
-------------------
It is a tool that helps you with the version control of several devtools. The name "gvc" stands for "go version control" or "general version control".

## Installation
-------------------
- Get executable file from github release or through "go install github.com/moqsien/gvc@latest".
- Then, execute the executable file, it will be automatically installed to "$HOME/.gvc".

## Features
-------------------
```bash
# gvc --help
gvc -h
```
```text
NAME:
   gvc - A new cli application

USAGE:
   gvc [global options] command [command options] [arguments...]

COMMANDS:
   host, h  gvc host
   go, g    gvc go <Command>
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help
```
```bash
# Get github dns, and modify hosts file.
# gvc host --help
gvc h -h
```
```text
moqsien@iMac-Pro gvc % gvc h -h
NAME:
   gvc host - gvc host

USAGE:
   gvc host [command options] [arguments...]

DESCRIPTION:
   Fetch hosts for github.

OPTIONS:
   --help, -h  show help
```
```bash
# Go version control.
# gvc go --help
gvc g -h
```
```text
NAME:
   gvc go - gvc go <Command>

USAGE:
   gvc go command [command options] [arguments...]

DESCRIPTION:
   Go version control.

COMMANDS:
   remote, r             gvc go r
   use, u                gvc go use
   local, l              gvc go local
   remove-unused, ru     gvc go ru
   remove-version, rm    gvc go rm
   add-envs, env, e, ae  gvc go env
   help, h               Shows a list of commands or help for one command

OPTIONS:
   --help, -h  show help
```
```bash
# Get specific version of go compilers from official website, install it, and set the envs for go.
# for example, gvc go use 1.20.1
# gvc go use --help
gvc g u -h
```
```text
NAME:
   gvc go use - gvc go use

USAGE:
   gvc go use [command options] [arguments...]

DESCRIPTION:
   Download and use version.

OPTIONS:
   --help, -h  show help
```

## Download
-------------------
- [MacOS](https://github.com/moqsien/gvc/releases/download/v0.1.0/macos-amd64.zip)
- [MacOS arm](https://github.com/moqsien/gvc/releases/download/v0.1.0/macos-arm64.zip)
- [Linux](https://github.com/moqsien/gvc/releases/download/v0.1.0/linux-amd64.zip)
- [Linux arm](https://github.com/moqsien/gvc/releases/download/v0.1.0/linux-arm64.zip)
- [Windows](https://github.com/moqsien/gvc/releases/download/v0.1.0/windows-amd64.zip)
- [Windows arm](https://github.com/moqsien/gvc/releases/download/v0.1.0/windows-arm64.zip)
