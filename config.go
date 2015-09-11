package cfg

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
)

// Config implements access to configuration values.
type Config struct {
	raw             []string
	comments        []string
	values          map[string]string
	keyValuePattern *regexp.Regexp
}

// NewConfig creates a new empty configuration.
func NewConfig() *Config {
	return &Config{
		raw:             make([]string, 0),
		comments:        make([]string, 0),
		values:          make(map[string]string),
		keyValuePattern: regexp.MustCompile(`\s*(\S+)\s*=\s*(\S+)\s*`),
	}
}

// NewConfigFromReader creates a new empty config and populates it
// with data parsed from the reader.
// If a error occurs when parsing the input an error is returned.
func NewConfigFromReader(r io.Reader) (*Config, error) {
	c := NewConfig()

	err := c.parse(r)
	if err != nil {
		return nil, err
	}

	return c, nil
}

// GetString returns the value for key as a string with new lines unescaped.
// If the key is not found an error is returned.
func (c *Config) GetString(key string) (string, error) {
	val, err := c.get(key)
	if err != nil {
		return "", err
	}

	uq := strings.Replace(val, "\\n", "\n", -1)
	return uq, nil
}

// GetInt returns the value for key as an int in decimal base.
// If the key is not found an error is returned.
// If the value can not be represented as an integer an error is returned.
func (c *Config) GetInt(key string) (int, error) {
	val, err := c.get(key)
	if err != nil {
		return 0, err
	}
	i, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("Invalid integer (%s)", err)
	}

	return int(i), nil
}

// GetFloat returns the value for key as a float64.
// If the key is not found an error is returned.
// If the value can not be represented as a float an error is returned.
func (c *Config) GetFloat(key string) (float64, error) {
	val, err := c.get(key)
	if err != nil {
		return 0, err
	}
	f, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return 0, fmt.Errorf("Invalid float (%s)", err)
	}

	return f, nil
}

// GetBool returns the value for key as a bool.
// If the key is not found an error is returned.
// If the value can not be represented as a boolean an error is returned.
func (c *Config) GetBool(key string) (bool, error) {
	val, err := c.get(key)
	if err != nil {
		return false, err
	}
	b, err := strconv.ParseBool(val)
	if err != nil {
		return false, fmt.Errorf("Invalid boolean (%s)", err)
	}

	return b, nil
}

// Comments returns all parsed comments in the config as a list of strings.
// Comments can not be modified programatically, but can be read.
// The comments are in the order they are defined in the source config.
func (c *Config) Comments() []string {
	// Should be a copy of the data as the user could modify the data,
	// but the comments are not modifiable so even if the user changes
	// the data it will not affect the configuration.
	return c.comments
}

// SetString creates or updates a value attached to key.
// Any new lines in the value are escaped.
func (c *Config) SetString(key, value string) {
	qval := strings.Replace(value, "\n", "\\n", -1)
	c.set(key, qval)
}

// SetInt creates or updates a value attached to key.
// All integer values are formated in decimal base.
func (c *Config) SetInt(key string, value int) {
	qval := strconv.FormatInt(int64(value), 10)
	c.set(key, qval)
}

// SetFloat creates or updates a value attached to key.
// The float value is formated without exponents eg. 3.14 not 3.14E+00.
func (c *Config) SetFloat(key string, value float64) {
	qval := strconv.FormatFloat(value, 'f', -1, 64)
	c.set(key, qval)
}

// SetBool creates or updates a value attached to key.
// The bool value is formated as "true" or "false".
func (c *Config) SetBool(key string, value bool) {
	qval := strconv.FormatBool(value)
	c.set(key, qval)
}

// Unset deletes a value from the config.
// Any comments defined in the source are preserved.
func (c *Config) Unset(key string) {
	for i, line := range c.raw {
		matches := c.keyValuePattern.FindStringSubmatch(line)
		if len(matches) == 3 {
			if matches[1] == key {
				c.raw = append(c.raw[:i], c.raw[i+1:]...)
				delete(c.values, key)
				return
			}
		}
	}
}

// String returns a string representation of the config.
// All comments and values are present.
// Whitespaces are preserved as they were in the source that were parsed if any.
func (c *Config) String() string {
	return strings.Join(c.raw, "\n")
}

// get is the internal getter that only operates on strings.
// Returns errors if the key is undefined.
func (c *Config) get(key string) (string, error) {
	if val, ok := c.values[key]; ok {
		return val, nil
	}

	return "", fmt.Errorf("No such key (%s)", key)
}

// set is the internal setter that only operates on strings.
// set must update both the cached map of values and the raw string data.
func (c *Config) set(key, value string) {
	if _, ok := c.values[key]; !ok { // If new value add it to the raw data
		c.raw = append(c.raw, fmt.Sprintf("%s = %s", key, value))
	} else { // If existing value update it
		for i, line := range c.raw {
			matches := c.keyValuePattern.FindStringSubmatch(line)
			if len(matches) == 3 {
				if matches[1] == key {
					c.raw[i] = strings.Replace(c.raw[i], matches[2], value, 1)
				}
			}
		}
	}
	c.values[key] = value // Update the cached value
}

// parse is the internal parser that extracts all values and comments from
// the input source.
// Returns error if the parsing fails.
func (c *Config) parse(r io.Reader) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		c.raw = append(c.raw, line)

		tline := strings.TrimSpace(line)
		if strings.HasPrefix(tline, "#") {
			comment := strings.TrimSpace(strings.TrimLeft(tline, "#"))
			c.comments = append(c.comments, comment)
		} else if strings.Contains(tline, "=") {
			parts := strings.SplitN(tline, "=", 2)
			c.values[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
