package cfg

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
)

// tagKey is used as the key for struct field tags
const tagKey = "cfg"

// Unmarshal parses the config data and stores the result in the
// value pointed to by v. v must be a pointer to a struct.
//
// Unmarshal matches incoming keys to either the struct field name or its
// tag, preferring an exact match but also accepting a case-insensitive match.
// Only exported fields can be populated. The tag value "-" is used to skip
// a field.
//
// If the type indicated in the struct field does not match the type in the
// config the field is skipped. Eg. the field type is int but contains a non
// numerical string value in the config data.
func Unmarshal(data []byte, v interface{}) error {
	// Parse the config
	buf := bytes.NewBuffer(data)
	c, err := NewConfigFromReader(buf)
	if err != nil {
		return fmt.Errorf("cfg: error parsing data %s", err)
	}

	return UnmarshalFromConfig(c, v)
}

// UnmarshalFromConfig strores the data in config in the value pointed to by v.
// v must be a pointer to a struct.
//
// UnmarshalFromConfig matches incoming keys to either the struct field name
// or its tag, preferring an exact match but also accepting
// a case-insensitive match.
// Only exported fields can be populated. The tag value "-" is used to skip
// a field.
//
// If the type indicated in the struct field does not match the type in the
// config the field is skipped. Eg. the field type is int but contains a non
// numerical string value in the config data.
func UnmarshalFromConfig(c *Config, v interface{}) error {
	// Check that the type v we will populate is a struct
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.Elem().Kind() != reflect.Struct {
		return errors.New("cfg: interface must be a pointer to struct")
	}

	// Dereference the pointer if it is one
	rv = rv.Elem()

	// Loop through all fields of the struct
	for i := 0; i < rv.NumField(); i++ {
		fv := rv.Field(i)        // Save the Value of the field
		sf := rv.Type().Field(i) // Save the StructField of the field

		// Check if the field should be skipped
		if sf.PkgPath != "" { // unexported
			continue
		}
		tag := sf.Tag.Get(tagKey)
		if tag == "-" {
			continue
		}

		// Loop through all keys and match them against the field
		// set the value if it matches.
		for key, _ := range c.values {
			// Check so the tag, or the name case insensitive matches, if not
			// go on to the next key
			if key != tag && bytes.EqualFold([]byte(key), []byte(sf.Name)) == false {
				continue
			}

			// Update the value with the correct type
			switch fv.Kind() {
			case reflect.Int:
				val, err := c.GetInt(key)
				if err == nil {
					fv.SetInt(int64(val))
				}
			case reflect.Float64:
				val, err := c.GetFloat(key)
				if err == nil {
					fv.SetFloat(val)
				}
			case reflect.Bool:
				val, err := c.GetBool(key)
				if err == nil {
					fv.SetBool(val)
				}
			case reflect.String:
				val, err := c.GetString(key)
				if err == nil {
					fv.SetString(val)
				}
			}
		}
	}

	return nil
}
