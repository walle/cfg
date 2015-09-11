package cfg_test

import (
	"bytes"
	"testing"

	"github.com/walle/cfg"
)

const configString = `
# This is a comment

# An integer value
answer = 42

# A float value
pi = 3.14
# A boolean value
is_active = true

# A string value
quotes = Alea iacta est\nEt tu, Brute?
`

const configStringWithEmpty = `Answer = 42
Pi = 0
is_active = false
quotes =
`

type MyConfig struct {
	Answer     int
	Pi         float64
	IsActive   bool   `cfg:"is_active"`
	Quotes     string `cfg:"quotes"`
	unexported string
	NotUsed    string `cfg:"-"`
}

func Test_Unmarshal(t *testing.T) {
	myConfig := &MyConfig{}
	err := cfg.Unmarshal([]byte(configString), myConfig)
	if err != nil {
		t.Errorf("Error unmarshaling data: %s\n", err)
	}

	if myConfig.Answer != 42 {
		t.Errorf("Expected %v, got %v\n", 42, myConfig.Answer)
	}

	if myConfig.Pi != 3.14 {
		t.Errorf("Expected %v, got %v\n", 3.14, myConfig.Pi)
	}

	if myConfig.IsActive != true {
		t.Errorf("Expected %v, got %v\n", true, myConfig.IsActive)
	}

	q := "Alea iacta est\nEt tu, Brute?"
	if myConfig.Quotes != q {
		t.Errorf("Expected %q, got %q\n", q, myConfig.Quotes)
	}
}

func Test_UnmarshalWithEmpty(t *testing.T) {
	myConfig := &MyConfig{}
	err := cfg.Unmarshal([]byte(configStringWithEmpty), myConfig)
	if err != nil {
		t.Errorf("Error unmarshaling data: %s\n", err)
	}

	if myConfig.Answer != 42 {
		t.Errorf("Expected %v, got %v\n", 42, myConfig.Answer)
	}

	if myConfig.Pi != 0 {
		t.Errorf("Expected %v, got %v\n", 0, myConfig.Pi)
	}

	if myConfig.IsActive != false {
		t.Errorf("Expected %v, got %v\n", false, myConfig.IsActive)
	}

	q := ""
	if myConfig.Quotes != q {
		t.Errorf("Expected %q, got %q\n", q, myConfig.Quotes)
	}
}

func Test_UnmarshalFromConfig(t *testing.T) {
	myConfig := &MyConfig{}
	config, err := cfg.NewConfigFromReader(bytes.NewBufferString(configString))
	if err != nil {
		t.Errorf("Error parsing the config: %s\n", err)
	}

	err = cfg.UnmarshalFromConfig(config, myConfig)
	if err != nil {
		t.Errorf("Error unmarshaling data: %s\n", err)
	}

	if myConfig.Answer != 42 {
		t.Errorf("Expected %v, got %v\n", 42, myConfig.Answer)
	}

	if myConfig.Pi != 3.14 {
		t.Errorf("Expected %v, got %v\n", 3.14, myConfig.Pi)
	}

	if myConfig.IsActive != true {
		t.Errorf("Expected %v, got %v\n", true, myConfig.IsActive)
	}

	q := "Alea iacta est\nEt tu, Brute?"
	if myConfig.Quotes != q {
		t.Errorf("Expected %q, got %q\n", q, myConfig.Quotes)
	}
}

func Test_UnmarshalToNotStruct(t *testing.T) {
	var i int
	err := cfg.Unmarshal([]byte(configString), i)
	if err == nil {
		t.Errorf("Did not get error when trying to unmarshal to int\n")
	}
}
