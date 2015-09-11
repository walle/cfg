// Package cfg implements a flat key value configuration that is read/writeable
// The format is simple and contains support for comments.
// The purpose is to be a simple format that is easily read/modified by both
// humans and computers. Comments are only modifiable by humans, but are readable
// for computers.
// The format is line based and all configurations are defined on their own line.
// The format has support for four different types of values.
// integers, floats, booleans and strings.
// Integers are defined in decimal base.
// Floats are defined without exponents e.g. 3.14 not 3.14E+00.
// Booleans are defined as the string representation "true" or "false".
// Strings are defined as is with all new lines escaped to only take up one line.
//
// Example configuration
//
// 		# This is a comment
//
//		# An integer value
//		answer = 42
//
//		# A float value
//		pi = 3.14
//
//		# A boolean value
//		is_active = true
//
//		# A string value
//		quotes = Alea iacta est\nEt tu, Brute?
//
// The package also contains functionality to encode/decode
// (marshal and unmarshal) data to structs defined in go.
// The object's default key string is the struct field name but can be specified
// in the struct field's tag value. The "cfg" key in the struct field's tag
// value is the key name. Use "-" to skip the field. Like in the encoding/json
// package.
package cfg
