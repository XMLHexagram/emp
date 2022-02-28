package emp

import (
	"fmt"
	"github.com/XMLHexagram/emp/empErr"
	"reflect"
	"strconv"
)

type Parser struct {
	config *Config
}

type Config struct {
	// ZeroFields, if set to true, will zero fields before writing them.
	// For example, a map will be emptied before decoded values are put in
	// it. If this is false, a map will be merged.
	ZeroFields bool

	// The tag name that emp reads for field names. This
	// defaults to "emp"
	TagName string

	Prefix string

	AllowEmpty bool

	DirectDefault bool

	ParseStringToArrayAndSlice func(s string) []string
}

func NewParser(config *Config) (*Parser, error) {
	if config == nil {
		config = &Config{}
	}

	if config.TagName == "" {
		config.TagName = "emp"
	}

	if config.ParseStringToArrayAndSlice == nil {
		config.ParseStringToArrayAndSlice = ParseStringToSlice
	}

	return &Parser{
		config: config,
	}, nil
}

func Parse(inputPtrInterface interface{}) error {
	config := &Config{}

	parser, err := NewParser(config)
	if err != nil {
		return err
	}

	return parser.Parse(inputPtrInterface)
}

func (self *Parser) Parse(StructPtrInterface interface{}) error {
	return self.parse(self.config.Prefix, "", "", self.config.DirectDefault, reflect.ValueOf(StructPtrInterface).Elem())
}

func (self *Parser) parse(prefix string, name string, default_ string, directDefault bool, outVal reflect.Value) error {
	var err error
	outValKind := getKind(outVal)
	switch outValKind {
	case reflect.Bool:
		err = self.parseBool(prefix, name, default_, directDefault, outVal)
	case reflect.Int:
		err = self.parseIntX(prefix, name, default_, directDefault, outVal, 0)
	case reflect.Int8:
		err = self.parseIntX(prefix, name, default_, directDefault, outVal, 8)
	case reflect.Int16:
		err = self.parseIntX(prefix, name, default_, directDefault, outVal, 16)
	case reflect.Int32:
		err = self.parseIntX(prefix, name, default_, directDefault, outVal, 32)
	case reflect.Int64:
		err = self.parseIntX(prefix, name, default_, directDefault, outVal, 64)
	case reflect.Uint:
		err = self.parseUintX(prefix, name, default_, directDefault, outVal, 0)
	case reflect.Uint8:
		err = self.parseUintX(prefix, name, default_, directDefault, outVal, 8)
	case reflect.Uint16:
		err = self.parseUintX(prefix, name, default_, directDefault, outVal, 16)
	case reflect.Uint32:
		err = self.parseUintX(prefix, name, default_, directDefault, outVal, 32)
	case reflect.Uint64:
		err = self.parseUintX(prefix, name, default_, directDefault, outVal, 64)
	case reflect.Float32:
		err = self.parseFloatX(prefix, name, default_, directDefault, outVal, 32)
	case reflect.Float64:
		err = self.parseFloatX(prefix, name, default_, directDefault, outVal, 64)
	case reflect.String:
		err = self.parseString(prefix, name, default_, directDefault, outVal)
	case reflect.Pointer:
		err = self.parsePointer(prefix, name, default_, directDefault, outVal)
	case reflect.Map:
		// TODO: wait for a great way
		err = self.parseMap(prefix, name, default_, outVal)
	case reflect.Struct:
		err = self.parseStruct(prefix, name, default_, directDefault, outVal)
	case reflect.Array:
		err = self.parseArray(prefix, name, default_, directDefault, outVal)
	case reflect.Slice:
		err = self.parseSlice(prefix, name, default_, directDefault, outVal)
	case reflect.Interface:
		err = self.parseInterface(prefix, name, default_, directDefault, outVal)
	}

	return err
}

func (self *Parser) parseBool(prefix string, name string, default_ string, directDefault bool, val reflect.Value) error {
	val = reflect.Indirect(val)
	// valType := val.Type()

	var value bool

	key := prefix + name
	envString, err := getEnvString(key, default_, directDefault, self.config.AllowEmpty)
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

func (self *Parser) parseString(prefix string, name string, default_ string, directDefault bool, val reflect.Value) error {
	val = reflect.Indirect(val)
	// valType := val.Type()

	var value string

	key := prefix + name
	envString, err := getEnvString(key, default_, directDefault, self.config.AllowEmpty)
	if err != nil {
		return err
	}

	value = envString

	val.SetString(value)
	return nil
}

func (self *Parser) parsePointer(prefix string, name string, default_ string, directDefault bool, val reflect.Value) error {
	// Create an element of the concrete (non pointer) type and decode
	// into that. Then set the value of the pointer to this type.
	valType := val.Type()
	valElemType := valType.Elem()

	if val.CanSet() {
		realVal := val
		if realVal.IsNil() || self.config.ZeroFields {
			realVal = reflect.New(valElemType)
		}

		err := self.parse(prefix, name, "", directDefault, reflect.Indirect(realVal))
		if err != nil {
			return err
		}

		val.Set(realVal)
	} else {
		err := self.parse(prefix, name, "", directDefault, reflect.Indirect(val))
		if err != nil {
			return err
		}
	}

	return nil
}

func (self *Parser) parseStruct(prefix string, name string, default_ string, directDefault bool, val reflect.Value) error {
	val = reflect.Indirect(val)
	valType := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldName := valType.Field(i).Name
		tagPrefix, name, default_, isIgnore := parseTagString(valType.Field(i).Tag.Get(self.config.TagName))

		if name == "" {
			name = fieldName
		}

		if isIgnore {
			continue
		}

		err := self.parse(prefix+tagPrefix, name, default_, directDefault, field)
		if err != nil {
			return err
		}
	}
	return nil
}

func (self *Parser) parseFloatX(prefix string, name string, default_ string, directDefault bool, val reflect.Value, X int) error {
	val = reflect.Indirect(val)
	// valType := val.Type()

	var value float64

	key := prefix + name
	envString, err := getEnvString(key, default_, directDefault, self.config.AllowEmpty)
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

func (self *Parser) parseIntX(prefix string, name string, default_ string, directDefault bool, val reflect.Value, X int) error {
	val = reflect.Indirect(val)
	// valType := val.Type()

	var value int64

	key := prefix + name
	envString, err := getEnvString(key, default_, directDefault, self.config.AllowEmpty)
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

func (self *Parser) parseUintX(prefix string, name string, default_ string, directDefault bool, val reflect.Value, X int) error {
	val = reflect.Indirect(val)
	// valType := val.Type()

	var value uint64

	key := prefix + name
	envString, err := getEnvString(key, default_, directDefault, self.config.AllowEmpty)
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

func (self *Parser) parseMap(prefix string, name string, default_ string, val reflect.Value) error {
	return empErr.UnsupportedTypeError.New().Wrap("map type is not supported")
}

func (self *Parser) parseArray(prefix string, name string, default_ string, directDefault bool, val reflect.Value) error {
	valType := val.Type()
	valElemType := valType.Elem()
	arrayType := reflect.ArrayOf(valType.Len(), valElemType)

	valArray := val

	key := prefix + name
	envString, err := getEnvString(key, default_, directDefault, self.config.AllowEmpty)
	if err != nil {
		return err
	}

	// Make a new array to hold our result, same size as the original data.
	if self.config.ZeroFields {
		valArray = reflect.New(arrayType).Elem()
	}
	dataSlice := self.config.ParseStringToArrayAndSlice(envString)

	if len(dataSlice) > valArray.Len() {
		return empErr.
			ArraySizeMismatchError.New().Wrap(fmt.Sprintf("'%s': expected source data to have length less or equal to %d, got %d", name, arrayType.Len(), len(dataSlice)))
	}

	// Accumulate any errors
	errors := make([]string, 0)

	for i, v := range dataSlice {
		err := self.parse("", "", v, true, valArray.Index(i))
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

func (self *Parser) parseSlice(prefix string, name string, default_ string, directDefault bool, val reflect.Value) error {
	valType := val.Type()
	valElemType := valType.Elem()
	sliceType := reflect.SliceOf(valElemType)

	valSlice := val
	if valSlice.IsNil() || self.config.ZeroFields {
		// Make a new slice to hold our result, same size as the original data.
		valSlice = reflect.MakeSlice(sliceType, 0, 0)
	}

	key := prefix + name
	envString, err := getEnvString(key, default_, directDefault, self.config.AllowEmpty)
	if err != nil {
		return err
	}

	dataSlice := self.config.ParseStringToArrayAndSlice(envString)

	// Accumulate any errors
	errors := make([]string, 0)

	for i, v := range dataSlice {
		for valSlice.Len() <= i {
			valSlice = reflect.Append(valSlice, reflect.Zero(valElemType))
		}
		currentField := valSlice.Index(i)

		err := self.parse("", "", v, true, currentField)
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

func (self *Parser) parseInterface(prefix string, name string, default_ string, directDefault bool, val reflect.Value) error {
	val = reflect.Indirect(val)
	// valType := val.Type()

	var value string

	key := prefix + name
	envString, err := getEnvString(key, default_, directDefault, self.config.AllowEmpty)
	if err != nil {
		return err
	}

	value = envString

	val.Set(reflect.ValueOf(value))
	return nil
}
