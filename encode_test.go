package cfg_test

import (
	"testing"

	"github.com/walle/cfg"
)

var myConfig = &MyConfig{
	Answer:     42,
	Pi:         3.14,
	IsActive:   true,
	Quotes:     "Alea iacta est\nEt tu, Brute?",
	unexported: "foo",
	NotUsed:    "bar",
}

var myConfigWithEmpty = &MyConfig{
	Answer:   42,
	IsActive: false,
}

const myConfigEncoded = `Answer = 42
Pi = 3.14
is_active = true
quotes = Alea iacta est\nEt tu, Brute?
`

const myConfigEncodedWithEmpty = `Answer = 42
Pi = 0
is_active = false
quotes = 
`

func Test_Marshal(t *testing.T) {
	data, err := cfg.Marshal(myConfig)
	if err != nil {
		t.Errorf("Error encoding data: %s\n", err)
	}

	if string(data) != myConfigEncoded {
		t.Errorf("Expected %q got %q\n", myConfigEncoded, string(data))
	}
}

func Test_MarshalWithEmpty(t *testing.T) {
	data, err := cfg.Marshal(myConfigWithEmpty)
	if err != nil {
		t.Errorf("Error encoding data: %s\n", err)
	}

	if string(data) != myConfigEncodedWithEmpty {
		t.Errorf("Expected %q got %q\n", myConfigEncodedWithEmpty, string(data))
	}
}

func Test_MarshalNotStruct(t *testing.T) {
	i := 0
	_, err := cfg.Marshal(i)
	if err == nil {
		t.Errorf("Did not get error when trying to marshal int\n")
	}
}

func Test_MarshalToConfig(t *testing.T) {
	config, err := cfg.MarshalToConfig(myConfig)
	if err != nil {
		t.Errorf("Error encoding data: %s\n", err)
	}

	answer, _ := config.GetInt("Answer")
	if answer != 42 {
		t.Errorf("Expected %v, got %v\n", 42, answer)
	}

	pi, _ := config.GetFloat("Pi")
	if pi != 3.14 {
		t.Errorf("Expected %v, got %v\n", 3.14, pi)
	}

	isActive, _ := config.GetBool("is_active")
	if isActive != true {
		t.Errorf("Expected %v, got %v\n", true, isActive)
	}

	quotes, _ := config.GetString("quotes")
	q := "Alea iacta est\nEt tu, Brute?"
	if quotes != q {
		t.Errorf("Expected %q, got %q\n", q, quotes)
	}
}
