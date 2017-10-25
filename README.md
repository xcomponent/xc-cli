# xc-cli
XComponent Command Line Interface

[![Coverage Status](https://coveralls.io/repos/github/xcomponent/xc-cli/badge.svg?branch=master)](https://coveralls.io/github/xcomponent/xc-cli?branch=master) [![Build Status](https://travis-ci.org/xcomponent/xc-cli.svg?branch=master)](https://travis-ci.org/xcomponent/xc-cli)

## Build

XC Cli requires [dep](https://github.com/golang/dep) to manage its dependencies.

```
$ go get -u github.com/golang/dep/cmd/dep
```

Once dep is available in your `PATH`, you are ready to build and run XC Cli

```
$ dep ensure
$ go install
```

```
$ xc-cli
NAME:
   XC CLI - XComponent Command Line Interface

USAGE:
   xc-cli [global options] command [command options] [arguments...]

VERSION:
   0.2.0

COMMANDS:
     install  Install XComponent
     init     Initialize a new XComponent project
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```

## Release

Install [ghr](https://github.com/tcnksm/ghr)

```
$ go get -u github.com/tcnksm/ghr
```

Generate a Github [Personal Access Token](https://github.com/settings/tokens/new) with the `public_repo` scope. Add that token to your git global configuration.

```
$ git config --global github.token "....."
```

You can now launch `RELEASE.sh` with the tag version you would like to release.

```
$ ./RELEASE.sh 0.2.0
```
