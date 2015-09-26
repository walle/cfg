package cfg_test

import (
	"bytes"
	"fmt"

	"github.com/walle/cfg"
)

func ExampleMarshal() {
	type MyConfig struct {
		Answer     int
		Pi         float64
		IsActive   bool   `cfg:"is_active"`
		Quotes     string `cfg:"quotes"`
		unexported string
		NotUsed    string `cfg:"-"`
	}

	var myConfig = &MyConfig{
		Answer:     42,
		Pi:         3.14,
		IsActive:   true,
		Quotes:     "Alea iacta est\nEt tu, Brute?",
		unexported: "foo",
		NotUsed:    "bar",
	}

	data, err := cfg.Marshal(myConfig)
	if err != nil {
		// Handle error
	}

	fmt.Printf("%s", data)
	// Output:
	// Answer = 42
	// Pi = 3.14
	// is_active = true
	// quotes = Alea iacta est\nEt tu, Brute?
}

func ExampleUnmarshal() {
	type MyConfig struct {
		Answer     int
		Pi         float64
		IsActive   bool   `cfg:"is_active"`
		Quotes     string `cfg:"quotes"`
		unexported string
		NotUsed    string `cfg:"-"`
	}

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

	myConfig := &MyConfig{}
	err := cfg.Unmarshal([]byte(configString), myConfig)
	if err != nil {
		// Handle error
	}

	fmt.Println(myConfig.Answer)
	fmt.Println(myConfig.Pi)
	fmt.Println(myConfig.IsActive)
	fmt.Println(myConfig.Quotes)

	// Output:
	// 42
	// 3.14
	// true
	// Alea iacta est
	// Et tu, Brute?
}

func ExampleNewConfigFromReader() {
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

	r := bytes.NewBufferString(configString)
	config, err := cfg.NewConfigFromReader(r)
	if err != nil {
		// Handle error
	}

	a, _ := config.GetInt("answer")
	fmt.Println(a)
	p, _ := config.GetFloat("pi")
	fmt.Println(p)
	i, _ := config.GetBool("is_active")
	fmt.Println(i)
	q, _ := config.GetString("quotes")
	fmt.Println(q)

	// Output:
	// 42
	// 3.14
	// true
	// Alea iacta est
	// Et tu, Brute?
}
