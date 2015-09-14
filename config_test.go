package cfg_test

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/walle/cfg"
)

func Test_Read(t *testing.T) {
	config := newConfigFromFile("read", t)

	foo, err := config.GetString("foo")
	if err != nil {
		t.Errorf("Key foo not found: %s\n", err)
	}
	if foo != "bar" {
		t.Errorf("Expected %s got %s\n", "bar", foo)
	}

	bar, err := config.GetString("bar")
	if err != nil {
		t.Errorf("Key bar not found: %s\n", err)
	}
	if bar != "foo" {
		t.Errorf("Expected %s got %s\n", "foo", bar)
	}
}

func Test_Types(t *testing.T) {
	config := newConfigFromFile("types", t)

	i, err := config.GetInt("integer")
	if err != nil {
		t.Errorf("Key integer not found: %s\n", err)
	}
	if i != 42 {
		t.Errorf("Expected %v got %v\n", 42, i)
	}

	f, err := config.GetFloat("float")
	if err != nil {
		t.Errorf("Key float not found: %s\n", err)
	}
	if f != 4.2 {
		t.Errorf("Expected %v got %v\n", 4.2, f)
	}

	b, err := config.GetBool("boolean")
	if err != nil {
		t.Errorf("Key boolean not found: %s\n", err)
	}
	if b != true {
		t.Errorf("Expected %v got %v\n", true, b)
	}

	s, err := config.GetString("string")
	if err != nil {
		t.Errorf("Key string not found: %s\n", err)
	}
	if s != "This is a string!" {
		t.Errorf("Expected %q got %q\n", "This is a string!", s)
	}
}

func Test_IncompatibleTypes(t *testing.T) {
	config := newConfigFromFile("types", t)

	_, err := config.GetInt("string")
	if err == nil {
		t.Errorf("Expected type error but got none")
	}

	_, err = config.GetFloat("string")
	if err == nil {
		t.Errorf("Expected type error but got none")
	}

	_, err = config.GetBool("string")
	if err == nil {
		t.Errorf("Expected type error but got none")
	}
}

func Test_UndefinedKey(t *testing.T) {
	config := newConfigFromFile("types", t)

	_, err := config.GetInt("undefined")
	if err == nil {
		t.Errorf("Expected not found error but got none")
	}

	_, err = config.GetFloat("undefined")
	if err == nil {
		t.Errorf("Expected not found error but got none")
	}

	_, err = config.GetBool("undefined")
	if err == nil {
		t.Errorf("Expected not found error but got none")
	}

	_, err = config.GetString("undefined")
	if err == nil {
		t.Errorf("Expected not found error but got none")
	}
}

func Test_Comments(t *testing.T) {
	config := newConfigFromFile("comments", t)

	if len(config.Comments()) != 1 {
		t.Errorf("Wrong number of comments\n")
	}

	if config.Comments()[0] != "Test data for verifying that comments work" {
		t.Errorf("Comment data is wrong\n")
	}
}

func Test_Format(t *testing.T) {
	config := newConfigFromFile("format", t)

	foo, err := config.GetString("foo")
	if err != nil {
		t.Errorf("Key foo not found: %s\n", err)
	}
	if foo != "bar" {
		t.Errorf("Expected %s got %s\n", "bar", foo)
	}

	foo2, err := config.GetString("foo2")
	if err != nil {
		t.Errorf("Key foo2 not found: %s\n", err)
	}
	if foo2 != "bar" {
		t.Errorf("Expected %s got %s\n", "bar", foo2)
	}

	foo3, err := config.GetString("foo3")
	if err != nil {
		t.Errorf("Key foo3 not found: %s\n", err)
	}
	if foo3 != "bar" {
		t.Errorf("Expected %s got %s\n", "bar", foo3)
	}
}

func Test_Whitespace(t *testing.T) {
	config := newConfigFromFile("whitespace", t)

	foo, err := config.GetString("foo")
	if err != nil {
		t.Errorf("Key foo not found: %s\n", err)
	}
	if foo != "bar" {
		t.Errorf("Expected %q got %q\n", "bar", foo)
	}

	bar, err := config.GetString("bar")
	if err != nil {
		t.Errorf("Key bar not found: %s\n", err)
	}
	if bar != "foo" {
		t.Errorf("Expected %q got %q\n", "foo", bar)
	}

	foobar, err := config.GetString("foobar")
	if err != nil {
		t.Errorf("Key foobar not found: %s\n", err)
	}
	if foobar != "baz" {
		t.Errorf("Expected %q got %q\n", "baz", foobar)
	}
}

func Test_Newlines(t *testing.T) {
	config := newConfigFromFile("newlines", t)

	golden := getGolden("newlines.txt", t)

	foo, err := config.GetString("foo")
	if err != nil {
		t.Errorf("Key foo not found: %s\n", err)
	}
	if foo != golden {
		t.Errorf("Expected %q got %q\n", golden, foo)
	}
}

func Test_Create(t *testing.T) {
	config := cfg.NewConfig()

	golden := getGolden("create.cfg", t)

	config.SetString("foo", "bar")

	if config.String() != golden {
		t.Errorf("Expected %q got %q\n", golden, config.String())
	}
}

func Test_Add(t *testing.T) {
	config := newConfigFromFile("add", t)

	golden := getGolden("add.cfg", t)

	config.SetString("bar", "foo")

	if config.String() != golden {
		t.Errorf("Expected %q got %q\n", golden, config.String())
	}
}

func Test_Update(t *testing.T) {
	config := newConfigFromFile("update", t)

	golden := getGolden("update.cfg", t)

	config.SetString("string", "foobar")
	config.SetInt("integer", 404)
	config.SetFloat("float", 3.14)
	config.SetBool("boolean", true)

	if config.String() != golden {
		t.Errorf("Expected %q got %q\n", golden, config.String())
	}
}

func Test_Delete(t *testing.T) {
	config := newConfigFromFile("delete", t)

	golden := getGolden("delete.cfg", t)

	config.Unset("bar")

	if config.String() != golden {
		t.Errorf("Expected %q got %q\n", golden, config.String())
	}
}

func Test_DeleteRetainsComments(t *testing.T) {
	config := newConfigFromFile("delete_retains_comments", t)

	golden := getGolden("delete_retains_comments.cfg", t)

	config.Unset("bar")

	if config.String() != golden {
		t.Errorf("Expected %q got %q\n", golden, config.String())
	}
}

func Test_WriteNewlines(t *testing.T) {
	config := cfg.NewConfig()

	b, err := ioutil.ReadFile("_testdata/write_newlines.txt")
	if err != nil {
		t.Errorf("Error opening test data: %s\n", err)
	}
	str := strings.TrimSpace(string(b))

	golden := getGolden("write_newlines.cfg", t)

	config.SetString("foo", str)

	if config.String() != golden {
		t.Errorf("Expected %q got %q\n", golden, config.String())
	}
}

func newConfigFromFile(filename string, t *testing.T) *cfg.Config {
	f, err := os.Open(fmt.Sprintf("_testdata/%s.cfg", filename))
	if err != nil {
		t.Errorf("Error reading test data: %s\n", err)
	}
	config, err := cfg.NewConfigFromReader(f)
	if err != nil {
		t.Errorf("Error creating config: %s\n", err)
	}

	return config
}

func getGolden(filename string, t *testing.T) string {
	b, err := ioutil.ReadFile(fmt.Sprintf("_testdata/%s.golden", filename))
	if err != nil {
		t.Errorf("Error reading golden file: %s\n", err)
	}
	return strings.TrimSpace(string(b))
}

// Test that errors in the reader is handled

// Simple implementation of io.Reader that only returns error on Read
type errorReader struct {
	err error
}

func (e errorReader) Read(p []byte) (int, error) {
	return 0, io.ErrClosedPipe
}

func Test_ParseError(t *testing.T) {
	_, err := cfg.NewConfigFromReader(errorReader{})
	if err == nil {
		t.Errorf("Expected parse error but got none")
	}
}
