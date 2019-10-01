# gowatch [![Build Status](https://travis-ci.org/msAlcantara/gowatch.svg?branch=master)](https://travis-ci.org/msalcantara/gowatch)

gowatch is a tool to watch for .go files changes and rebuild automaticaly

## Installation

```bash
$ go get -u github.com/msalcantara/gowatch/cmd/gowatch
```

## Simple Usage
In you project path just type gowatch

```bash
$ gowatch
```

## Watch with custon flags
  - #### Flags to your binary command:
   - `$ gowatch apparg1 apparg2`

- #### Flags to go build command:
   - `$ gowatch --build-flags -build-flag1,build-flag2`

 - #### Watch in specific directory(default is current):
   - `$ gowatch -d ./custon/path`


## License
[MIT](https://github.com/msAlcantara/gowatch/blob/master/LICENSE)
