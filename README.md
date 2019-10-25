# gowatch [![Build Status](https://travis-ci.org/msAlcantara/gowatch.svg?branch=master)](https://travis-ci.org/msalcantara/gowatch) [![Coverage Status](https://coveralls.io/repos/github/msAlcantara/gowatch/badge.svg?branch=master)](https://coveralls.io/github/msAlcantara/gowatch?branch=master)

gowatch is a tool to watch for .go files changes and rebuild automaticaly

## Installation

```bash
$ go get -u github.com/msalcantara/gowatch/cmd/gowatch
```

## Usage
In you project path just type gowatch

```
$ gowatch
```

For help use `-h`or `--help`

```
$ gowatch -h
```

For version

```
$ gowatch --version
```

Pass arguments to your app

```
$ gowatch apparg1 apparg2
```

Or

```
$ gowatch --run-flags="-c file_config.conf"
```

Use custon args to `go build` command

```
$ gowatch --build-flags=-x,-v
```

Watch custon directory (default is current)

```
$ gowatch -d ./custon/path
```

gowatch restart in any .go files changes. To ignore some pattern of files use:

```
$ gowatch -i *_test.go
```

To show debug info of gowatch

```
$ gowatch -V
```

## Default Config

To use an default config, create a file called `.gowatch.yml` or use `--config` flag to pass a different file name. CLI flags override values of the file.

This file can have any flag that CLI accept.

```yaml
verbose: true


dir: .

ignore:
  - "*_test.go"

build_flags:
  - -x
  - -v


run_flags:
  - localhost
  - 8000

```


## License
[MIT](https://github.com/msAlcantara/gowatch/blob/master/LICENSE)
