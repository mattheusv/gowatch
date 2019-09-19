# gowatch

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
 - #### Flags to go build command:
   - `$ gowatch --build-flags -build-flag1,build-flag2`

 - #### Flags to your binary command:
   - `$ gowatch --run-flags your-custon-flag1,your-custon-flag2`

 - #### Watch in specific directory(default is current):
   - `$ gowatch -d /custon/path`


## License
[MIT](https://github.com/msAlcantara/gowatch/blob/master/LICENSE)
