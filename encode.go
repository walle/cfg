package cfg

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// Marshal returns the config encoding of v.
// v must be a pointer to a struct.
//
// Struct values encode as config values. Each exported struct field
// becomes a member of the object unless
//   - the field's tag is "-"
//
// Only top level values of the supported types are encoded. No recursion
// down the struct is done.
//
// The object's default key string is the struct field name
// but can be specified in the struct field's tag value. The "cfg" key in
// the struct field's tag value is the key name.
// Examples:
//
//   // Field is ignored by this package.
//   Field int `cfg:"-"`
//
//   // Field appears in config as key "myName".
//   Field int `cfg:"myName"`
func Marshal(v interface{}) ([]byte, error) {
	// Check that the type v we will read is a struct
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.Elem().Kind() != reflect.Struct {
		return []byte{}, errors.New("cfg: interface must be a pointer to struct")
	}

	// Dereference the pointer
	rv = rv.Elem()

	buf := bytes.NewBuffer([]byte{})

	// Loop through all fields of the struct
	for i := 0; i < rv.NumField(); i++ {
		fv := rv.Field(i)
		sf := rv.Type().Field(i)

		// Check if the field should be skipped
		if sf.PkgPath != "" { // unexported
			continue
		}
		tag := sf.Tag.Get(tagKey)
		if tag == "-" {
			continue
		}

		key := sf.Name
		if tag != "" {
			key = tag
		}

		err := writeValue(buf, &fv, key)
		if err != nil {
			return nil, fmt.Errorf("cfg: error writing value: %s", err)
		}
	}
	return buf.Bytes(), nil
}

// MarshalToConfig creates a config object from the marshaled data in v.
// v must be a pointer to a struct.
//
// Struct values encode as config values. Each exported struct field
// becomes a member of the object unless
//   - the field's tag is "-"
//
// Only top level values of the supported types are encoded. No recursion
// down the struct is done.
//
// The object's default key string is the struct field name
// but can be specified in the struct field's tag value. The "cfg" key in
// the struct field's tag value is the key name.
// Examples:
//
//   // Field is ignored by this package.
//   Field int `cfg:"-"`
//
//   // Field appears in config as key "myName".
//   Field int `cfg:"myName"`
func MarshalToConfig(v interface{}) (*Config, error) {
	data, err := Marshal(v)
	if err != nil {
		return nil, err
	}

	c, err := NewConfigFromReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	return c, nil
}

// writeValue adds the key value to buffer if it is exported and not skipped.
func writeValue(buf *bytes.Buffer, fv *reflect.Value, key string) error {
	switch fv.Kind() {
	case reflect.Int:
		_, err := buf.WriteString(fmt.Sprintf("%s = %v\n", key, fv.Int()))
		return err
	case reflect.Float64:
		_, err := buf.WriteString(fmt.Sprintf("%s = %v\n", key, fv.Float()))
		return err
	case reflect.Bool:
		_, err := buf.WriteString(fmt.Sprintf("%s = %v\n", key, fv.Bool()))
		return err
	case reflect.String:
		qval := strings.Replace(fv.String(), "\n", "\\n", -1)
		_, err := buf.WriteString(fmt.Sprintf("%s = %s\n", key, qval))
		return err
	}

	return nil
}
