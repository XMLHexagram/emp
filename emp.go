// Package emp exposes functionality to convert environment value into
// a struct.
//
// The Go structure can be arbitrarily complex, containing slices,
// other structs, etc. and the parser will properly parse the environment
// value and populate the values into the native struct.
//
// See the examples to see what the parser is capable of.
//
// The simplest function to start with is Parser.
//
// Field Tags
//
// When decoding to a struct, emp will use the field name by
// default to perform the mapping. For example, if a struct has a field
// "DATABASE_URL" then emp will look for a key in the environment
// of "DATABASE_URL" (case-sensitive).
//
//     type Model struct {
//         DATABASE_URL string
//     }
//
// You can change the behavior of emp by using struct tags.
// The default struct tag that emp looks for is "emp"
// but you can customize it using Config.
//
// Renaming Fields
//
// To rename the key that emp looks for, use the "emp"
// tag and set a value directly. For example, to change the "DATABASE_URL" example
// above to "DATABASE_DSN":
//
//     type Model struct {
//         DATABASE_URL string `emp:"DATABASE_DSN"`
//     }
//
// Embedded Structs and Squashing
//
// By default, Embedded structs are treated as prefix by their name.
// emp:
//
//    type Model struct {
//        DATABASE_URL string
//        USER_ User
//    }
//
//    type User struct {
//        Name string
//		  Pass string
//	  }
//
// With environment value:
//
//	DATABASE_URL=postgres://user:pass@host:port/db
//	USER_NAME=user
//	USER_PASS=pass
//
// Will be converted to:
//
//    type Model struct {
//        DATABASE_URL: "postgres://user:pass@host:port/db",
//        USER_: {
//            Name: "user",
//            Pass: "pass",
//        },
//    }
//
// Config has a field that changes the behavior of emp
// to disable auto prefix.
//
// Unexported fields
//
// Since unexported (private) struct fields cannot be set outside the package
// where they are defined, the parser will simply skip them.
//
// For this output type definition:
//
//     type Exported struct {
//         private string // this unexported field will be skipped
//         PUBLIC string
//     }
//
// Using this as environment value:
//
//     private=secret
//     PUBLIC=SECRET
//
// The following struct will be parsed:
//
//     type Exported struct {
//         private: "" // field is left with an empty string (zero value)
//         Public: "SECRET"
//     }
//
// Other Configuration
//
// emp is highly configurable. See the Config struct
// for other features and options that are supported.
package emp

import (
	"fmt"
	"github.com/XMLHexagram/emp/empErr"
	"reflect"
	"strconv"
)

// A Parser takes a raw interface value and fill it data,
// keeping track of rich error information along the way in case
// anything goes wrong. You can more finely control how the Parser
// behaves using the Config structure. The top-level parse
// method is just a convenience that sets up the most basic Parser.
type Parser struct {
	config *Config
}

// Config is the configuration that is used to create a new parser
// and allows customization of various aspects of decoding.
type Config struct {
	// ZeroFields, if set to true, will zero fields before writing them.
	// For example, a map will be emptied before parsed values are put in
	// it. If this is false, a map will be merged.
	ZeroFields bool

	// The TagName that emp reads for field names. This defaults to "emp".
	TagName string

	// Prefix, if set, will be the first part of all environment value name.
	Prefix string

	// AutoPrefix, default to true, it will add prefix automatically when meet embedded struct(use key).
	AutoPrefix bool

	// AllowEmpty, if set to true, will allow empty values in environment values.
	AllowEmpty bool

	// DirectDefault, if set to true, will use the default value in field name directly.
	DirectDefault bool

	// ParseStringToArrayAndSlice, customize the way split string to array and slice.
	ParseStringToArrayAndSlice func(s string) []string

	marshal    bool
	marshalRes string
}

// NewParser returns a new parser for the given configuration. Once
// a parser has been returned, the same configuration must not be used
// again.
func NewParser(config *Config) (*Parser, error) {
	if config == nil {
		config = &Config{}
	}

	if config.TagName == "" {
		config.TagName = "emp"
	}

	if config.ParseStringToArrayAndSlice == nil {
		config.ParseStringToArrayAndSlice = ParseStringToArrayAndSlice
	}

	return &Parser{
		config: config,
	}, nil
}

// Parse takes an input structure and uses reflection to translate it to
// the output structure. output must be a pointer to a struct.
func Parse(inputPtrInterface interface{}) error {
	config := &Config{}

	parser, err := NewParser(config)
	if err != nil {
		return err
	}

	return parser.Parse(inputPtrInterface)
}

// Marshal struct to get an env file format string.
func Marshal(inputPtrInterface interface{}) (string, error) {
	config := &Config{}

	parser, err := NewParser(config)
	if err != nil {
		return "", err
	}

	return parser.Marshal(inputPtrInterface)
}

// Parse parses the given raw interface to the target pointer specified
// by the configuration.
func (p *Parser) Parse(StructPtrInterface interface{}) error {
	return p.parse(p.config.Prefix, "", "", p.config.DirectDefault, reflect.ValueOf(StructPtrInterface).Elem())
}

// Marshal struct to get an env file format string.
func (p *Parser) Marshal(StructPtrInterface interface{}) (string, error) {
	p.config.marshal = true
	p.config.marshalRes = ""
	defer func() {
		p.config.marshal = false
		p.config.marshalRes = ""
	}()
	err := p.parse(p.config.Prefix, "", "", p.config.DirectDefault, reflect.ValueOf(StructPtrInterface).Elem())
	if err != nil {
		return "", err
	}
	return p.config.marshalRes, nil
}

// parse environment value to specific reflection value.
func (p *Parser) parse(prefix string, name string, default_ string, directDefault bool, outVal reflect.Value) error {
	var err error
	outValKind := getKind(outVal)
	switch outValKind {
	case reflect.Bool:
		err = p.parseBool(prefix, name, default_, directDefault, outVal)
	case reflect.Int:
		err = p.parseIntX(prefix, name, default_, directDefault, outVal, 0)
	case reflect.Int8:
		err = p.parseIntX(prefix, name, default_, directDefault, outVal, 8)
	case reflect.Int16:
		err = p.parseIntX(prefix, name, default_, directDefault, outVal, 16)
	case reflect.Int32:
		err = p.parseIntX(prefix, name, default_, directDefault, outVal, 32)
	case reflect.Int64:
		err = p.parseIntX(prefix, name, default_, directDefault, outVal, 64)
	case reflect.Uint:
		err = p.parseUintX(prefix, name, default_, directDefault, outVal, 0)
	case reflect.Uint8:
		err = p.parseUintX(prefix, name, default_, directDefault, outVal, 8)
	case reflect.Uint16:
		err = p.parseUintX(prefix, name, default_, directDefault, outVal, 16)
	case reflect.Uint32:
		err = p.parseUintX(prefix, name, default_, directDefault, outVal, 32)
	case reflect.Uint64:
		err = p.parseUintX(prefix, name, default_, directDefault, outVal, 64)
	case reflect.Float32:
		err = p.parseFloatX(prefix, name, default_, directDefault, outVal, 32)
	case reflect.Float64:
		err = p.parseFloatX(prefix, name, default_, directDefault, outVal, 64)
	case reflect.String:
		err = p.parseString(prefix, name, default_, directDefault, outVal)
	case reflect.Ptr:
		err = p.parsePointer(prefix, name, default_, directDefault, outVal)
	case reflect.Map:
		// TODO: wait for a great way
		err = p.parseMap(prefix, name, default_, outVal)
	case reflect.Struct:
		err = p.parseStruct(prefix, name, default_, directDefault, outVal)
	case reflect.Array:
		err = p.parseArray(prefix, name, default_, directDefault, outVal)
	case reflect.Slice:
		err = p.parseSlice(prefix, name, default_, directDefault, outVal)
	case reflect.Interface:
		err = p.parseInterface(prefix, name, default_, directDefault, outVal)
	}

	return err
}

func (p *Parser) parseBool(prefix string, name string, default_ string, directDefault bool, val reflect.Value) error {
	val = reflect.Indirect(val)
	// valType := val.Type()

	if p.config.marshal {
		p.config.marshalRes += fmt.Sprintf("%s=%t\n", prefix+name, val.Bool())
		return nil
	}

	var value bool

	key := prefix + name
	envString, err := getEnvString(key, default_, directDefault, p.config.AllowEmpty)
	if err != nil {
		return err
	}

	value, err = strconv.ParseBool(envString)
	if err != nil {
		return empErr.CannotParseEnvStringToTypeError.New().Wrap(err)
	}

	val.SetBool(value)
	return nil
}

func (p *Parser) parseString(prefix string, name string, default_ string, directDefault bool, val reflect.Value) error {
	val = reflect.Indirect(val)
	// valType := val.Type()

	if p.config.marshal {
		p.config.marshalRes += fmt.Sprintf("%s=%s\n", prefix+name, val.String())
		return nil
	}

	var value string

	key := prefix + name
	envString, err := getEnvString(key, default_, directDefault, p.config.AllowEmpty)
	if err != nil {
		return err
	}

	value = envString

	val.SetString(value)
	return nil
}

func (p *Parser) parsePointer(prefix string, name string, default_ string, directDefault bool, val reflect.Value) error {
	// Create an element of the concrete (non pointer) type and decode
	// into that. Then set the value of the pointer to this type.
	valType := val.Type()
	valElemType := valType.Elem()

	if p.config.marshal {
		err := p.parse(prefix, name, "", directDefault, reflect.Indirect(val))
		return err
	}

	if val.CanSet() {
		realVal := val
		if realVal.IsNil() || p.config.ZeroFields {
			realVal = reflect.New(valElemType)
		}

		err := p.parse(prefix, name, "", directDefault, reflect.Indirect(realVal))
		if err != nil {
			return err
		}

		val.Set(realVal)
	} else {
		err := p.parse(prefix, name, "", directDefault, reflect.Indirect(val))
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *Parser) parseStruct(prefix string, name string, default_ string, directDefault bool, val reflect.Value) error {
	val = reflect.Indirect(val)
	valType := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		if !field.CanSet() {
			continue
		}
		fieldName := valType.Field(i).Name
		tagPrefix, name, default_, isIgnore := parseTagString(valType.Field(i).Tag.Get(p.config.TagName))

		if name == "" {
			name = fieldName
		}

		if isIgnore {
			continue
		}

		// auto prefix
		if field.Type().Kind() == reflect.Struct && tagPrefix == "" && p.config.AutoPrefix {
			prefix = name
		}

		err := p.parse(prefix+tagPrefix, name, default_, directDefault, field)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Parser) parseFloatX(prefix string, name string, default_ string, directDefault bool, val reflect.Value, X int) error {
	val = reflect.Indirect(val)
	// valType := val.Type()

	if p.config.marshal {
		p.config.marshalRes += fmt.Sprintf("%s=%f", prefix+name, val.Float())
		return nil
	}

	var value float64

	key := prefix + name
	envString, err := getEnvString(key, default_, directDefault, p.config.AllowEmpty)
	if err != nil {
		return err
	}

	value, err = strconv.ParseFloat(envString, X)
	if err != nil {
		return empErr.CannotParseEnvStringToTypeError.New().Wrap(err)
	}

	val.SetFloat(value)
	return nil
}

func (p *Parser) parseIntX(prefix string, name string, default_ string, directDefault bool, val reflect.Value, X int) error {
	val = reflect.Indirect(val)
	// valType := val.Type()

	if p.config.marshal {
		p.config.marshalRes += fmt.Sprintf("%s=%d\n", prefix+name, val.Int())
		return nil
	}

	var value int64

	key := prefix + name
	envString, err := getEnvString(key, default_, directDefault, p.config.AllowEmpty)
	if err != nil {
		return err
	}

	value, err = strconv.ParseInt(envString, 0, X)
	if err != nil {
		return empErr.CannotParseEnvStringToTypeError.New().Wrap(err)
	}

	val.SetInt(value)
	return nil
}

func (p *Parser) parseUintX(prefix string, name string, default_ string, directDefault bool, val reflect.Value, X int) error {
	val = reflect.Indirect(val)
	// valType := val.Type()

	if p.config.marshal {
		p.config.marshalRes += fmt.Sprintf("%s=%d\n", prefix+name, val.Uint())
		return nil
	}

	var value uint64

	key := prefix + name
	envString, err := getEnvString(key, default_, directDefault, p.config.AllowEmpty)
	if err != nil {
		return err
	}

	value, err = strconv.ParseUint(envString, 0, X)
	if err != nil {
		return empErr.CannotParseEnvStringToTypeError.New().Wrap(err)
	}

	val.SetUint(value)
	return nil
}

func (p *Parser) parseMap(prefix string, name string, default_ string, val reflect.Value) error {
	return empErr.UnsupportedTypeError.New().Wrap("map type is not supported")
}

func (p *Parser) parseArray(prefix string, name string, default_ string, directDefault bool, val reflect.Value) error {
	valType := val.Type()
	valElemType := valType.Elem()
	arrayType := reflect.ArrayOf(valType.Len(), valElemType)

	if p.config.marshal {
		p.config.marshalRes += fmt.Sprintf("%s=%v\n", prefix+name, formatSliceAndArrayReflectValue(val))
		return nil
	}

	valArray := val

	key := prefix + name
	envString, err := getEnvString(key, default_, directDefault, p.config.AllowEmpty)
	if err != nil {
		return err
	}

	// Make a new array to hold our result, same size as the original data.
	if p.config.ZeroFields {
		valArray = reflect.New(arrayType).Elem()
	}
	dataSlice := p.config.ParseStringToArrayAndSlice(envString)

	if len(dataSlice) > valArray.Len() {
		return empErr.
			ArraySizeMismatchError.New().Wrap(fmt.Sprintf("'%s': expected source data to have length less or equal to %d, got %d", name, arrayType.Len(), len(dataSlice)))
	}

	// Accumulate any errors
	errors := make([]string, 0)

	for i, v := range dataSlice {
		err := p.parse("", "", v, true, valArray.Index(i))
		if err != nil {
			errors = append(errors, err.Error())
		}
	}

	val.Set(valArray)

	if len(errors) > 0 {
		return empErr.CannotParseEnvStringToTypeError.New().Wrap(errors)
	}
	return nil
}

func (p *Parser) parseSlice(prefix string, name string, default_ string, directDefault bool, val reflect.Value) error {
	valType := val.Type()
	valElemType := valType.Elem()
	sliceType := reflect.SliceOf(valElemType)

	if p.config.marshal {
		p.config.marshalRes += fmt.Sprintf("%s=%v\n", prefix+name, formatSliceAndArrayReflectValue(val))
		return nil
	}

	valSlice := val
	if valSlice.IsNil() || p.config.ZeroFields {
		// Make a new slice to hold our result, same size as the original data.
		valSlice = reflect.MakeSlice(sliceType, 0, 0)
	}

	key := prefix + name
	envString, err := getEnvString(key, default_, directDefault, p.config.AllowEmpty)
	if err != nil {
		return err
	}

	dataSlice := p.config.ParseStringToArrayAndSlice(envString)

	// Accumulate any errors
	errors := make([]string, 0)

	for i, v := range dataSlice {
		for valSlice.Len() <= i {
			valSlice = reflect.Append(valSlice, reflect.Zero(valElemType))
		}
		currentField := valSlice.Index(i)

		err := p.parse("", "", v, true, currentField)
		if err != nil {
			errors = append(errors, err.Error())
		}
	}

	val.Set(valSlice)

	if len(errors) > 0 {
		return empErr.CannotParseEnvStringToTypeError.New().Wrap(errors)
	}
	return nil
}

func (p *Parser) parseInterface(prefix string, name string, default_ string, directDefault bool, val reflect.Value) error {
	val = reflect.Indirect(val)
	// valType := val.Type()

	if p.config.marshal {
		p.config.marshalRes += fmt.Sprintf("%s=%v\n", prefix+name, val.Interface())
		return nil
	}

	var value string

	key := prefix + name
	envString, err := getEnvString(key, default_, directDefault, p.config.AllowEmpty)
	if err != nil {
		return err
	}

	value = envString

	val.Set(reflect.ValueOf(value))
	return nil
}
