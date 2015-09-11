package cfg

import (
	"fmt"
	"os"
)

// ConfigFile is a utility type that can load and save config to a file.
type ConfigFile struct {
	path string
	*Config
}

// NewConfigFile returns a new ConfigFile with the parsed data in
// the file at path. Returns an error if the file can't be read or
// if the parsing of the config fails.
func NewConfigFile(path string) (*ConfigFile, error) {
	f, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("Could not open file: %s", err)
	}
	c, err := NewConfigFromReader(f)
	if err != nil {
		return nil, fmt.Errorf("Could not parse file: %s", err)
	}
	f.Close()

	return &ConfigFile{path: path, Config: c}, nil
}

// Persist saves all configured values to the file.
// Returns error if something goes wrong.
func (c *ConfigFile) Persist() error {
	f, err := os.OpenFile(c.path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("Could not open file: %s", err)
	}
	f.WriteString(c.String())
	f.Close()

	return nil
}

// Path returns the path to the file with the config.
func (c *ConfigFile) Path() string {
	return c.path
}
