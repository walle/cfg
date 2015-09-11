[![Build Status](https://img.shields.io/travis/walle/cfg.svg?style=flat)](https://travis-ci.org/walle/cfg)
[![Coverage](https://img.shields.io/codecov/c/github/walle/cfg.svg?style=flat)](https://codecov.io/github/walle/cfg)
[![Godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/walle/cfg)
[![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/walle/cfg/master/LICENSE)

# cfg - Simple read/write and commentable config files

Package cfg implements a flat key value configuration that is read/writeable.

The format is simple and contains support for comments. The purpose is to be a
simple format that is easily read/modified by both humans and computers.

Comments are only modifiable by humans, but are readable for computers. 

The format is line based and all configurations are defined on their own line. 

The format has support for four different types of values. integers, floats,
booleans and strings. Integers are defined in decimal base. Floats are defined
without exponents e.g. 3.14 not 3.14E+00. Booleans are defined as the string
representation "true" or "false". Strings are defined as is with all new lines
escaped to only take up one line.

Example configuration

```
# This is a comment

# An integer value
answer = 42

# A float value
pi = 3.14

# A boolean value
is_active = true

# A string value
quotes = Alea iacta est\nEt tu, Brute?
```

The package also contains functionality to encode/decode (marshal and
unmarshal) data to structs defined in go. The object's default key string is
the struct field name but can be specified in the struct field's tag value.
The "cfg" key in the struct field's tag value is the key name. Use "-" to skip
the field. Like in the encoding/json package.

## Why?

There are many good libraries for reading configurations, but none (that I know of) that allows the user to update values in the config file.
Go have really good support for JSON and JSON-files can be used for
configuration. But they have a drawback, and that is that they cannot be
commented (easily).

This library fixes this by having config files with comment support that also
can be updated from code.

## The format

The config file format is line based. Any line can be empty, contain a comment
or contain a key value pair.

* A comment is any line that starts with the character `#` excluding whitespace
* A configuration option is any line in the format key = value
* All other lines are empty

The parser skips all whitespace, so whitespace before or after data is discarded.
The key value, or configuration options can have as many whitespace characters
before, after or in between the key value as desired. The following examples
are all valid, and the foo[2-6]? keys all contain just "foo bar"

```
foo=foo bar
foo2= foo bar
foo3 =foo bar
foo4           =                   foo bar
foo5 = foo bar
    foo6 = foo bar
```

String values that contain line breaks have the line breaks escaped with the
`\n` character. The config handles the escaping and unescapes the values when
they are accessed.

## Installation

To install cfg, just use `go get`.

```shell
$ go get github.com/walle/cfg
```

To start using it import the package.

```go
import "github.com/walle/cfg"
```

## Examples

### Config example

```
# This is a comment

# An integer value
answer = 42

# A float value
pi = 3.14

# A boolean value
is_active = true

# A string value
quotes = Alea iacta est\nEt tu, Brute?
```

### Usage example

There are more examples defined in the [example_test.go](example_test.go)
file. And you can easily view them in the documentation
http://godoc.org/github.com/walle/cfg#pkg-examples

```go
myConfig := "foo = bar"
r := bytes.NewBufferString(myConfig)
config, err := cfg.NewConfigFromReader(r)
if err != nil {
        // Could not parse config
}
foo, err := config.GetString("foo")
if err != nil {
        // Could not find the key foo
}
fmt.Println(foo)
// Output:
// bar
```

## Testing

Run the tests using `go test`.

```shell
$ go test -cover
```

## Contributing

All contributions are welcome! See [CONTRIBUTING](CONTRIBUTING.md) for more
info.

## License

The code is under the MIT license. See [LICENSE](LICENSE) for more
information.

## Authors

The code is written by Fredrik Wallgren - https://github.com/walle
