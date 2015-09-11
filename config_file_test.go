package cfg_test

import (
	"io/ioutil"
	"testing"

	"github.com/walle/cfg"
)

const configContents = `
# This is a comment

# An integer value
answer = 42

# A float value
pi = 3.14
# A boolean value
is_active = true

# A string value
quotes = Alea iacta est\nEt tu, Brute?`

func Test_ConfigFile(t *testing.T) {
	f, err := ioutil.TempFile("", "cfg-test")
	if err != nil {
		t.Errorf("Error creating tmp file: %s\n", err)
	}

	path := f.Name()
	f.WriteString(configContents)
	f.Close()

	configFile, err := cfg.NewConfigFile(path)
	if err != nil {
		t.Errorf("Error parsing config: %s\n", err)
	}

	a, _ := configFile.GetInt("answer")
	if a != 42 {
		t.Errorf("Expected %v got %v\n", 42, a)
	}

	p, _ := configFile.GetFloat("pi")
	if p != 3.14 {
		t.Errorf("Expected %v got %v\n", 3.14, p)
	}

	i, _ := configFile.GetBool("is_active")
	if i != true {
		t.Errorf("Expected %v got %v\n", true, i)
	}

	quotes := "Alea iacta est\nEt tu, Brute?"
	q, _ := configFile.GetString("quotes")
	if q != quotes {
		t.Errorf("Expected %s got %s\n", quotes, q)
	}
}

func Test_ConfigFilePersist(t *testing.T) {
	f, err := ioutil.TempFile("", "cfg-test")
	if err != nil {
		t.Errorf("Error creating tmp file: %s\n", err)
	}

	path := f.Name()
	f.WriteString(configContents)
	f.Close()

	configFile, err := cfg.NewConfigFile(path)
	if err != nil {
		t.Errorf("Error parsing config: %s\n", err)
	}

	a, _ := configFile.GetInt("answer")
	if a != 42 {
		t.Errorf("Expected %v got %v\n", 42, a)
	}

	configFile.SetInt("answer", 314)

	err = configFile.Persist()
	if err != nil {
		t.Errorf("Error persisting config: %s\n", err)
	}

	c2, err := cfg.NewConfigFile(path)
	if err != nil {
		t.Errorf("Error parsing config: %s\n", err)
	}

	a2, _ := c2.GetInt("answer")
	if a2 != 314 {
		t.Errorf("Expected %v got %v\n", 314, a2)
	}
}

func Test_ConfigFilePath(t *testing.T) {
	f, err := ioutil.TempFile("", "cfg-test")
	if err != nil {
		t.Errorf("Error creating tmp file: %s\n", err)
	}

	path := f.Name()
	f.WriteString(configContents)
	f.Close()

	configFile, err := cfg.NewConfigFile(path)
	if err != nil {
		t.Errorf("Error parsing config: %s\n", err)
	}

	if configFile.Path() != path {
		t.Errorf("Expected %s, got %s\n", path, configFile.Path())
	}
}
