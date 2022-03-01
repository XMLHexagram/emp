package emp

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"strings"
	"testing"
)

func parseEnv(envMap map[string]string) {
	for k, v := range envMap {
		err := os.Setenv(k, v)
		if err != nil {
			panic(err)
		}
	}
}

func TestBool(t *testing.T) {
	parseEnv(map[string]string{
		"TEST_BOOL_1": "true",
		"TEST_BOOL_2": "false",
	})

	type args struct {
		TEST_BOOL_1 bool
		TEST_BOOL_2 bool
	}

	expect := &args{
		TEST_BOOL_1: true,
		TEST_BOOL_2: false,
	}

	res := new(args)

	err := Parse(res)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expect, res)
}

func TestIntX(t *testing.T) {
	parseEnv(map[string]string{
		"TEST_INT":   "114514",
		"TEST_INT8":  "127",
		"TEST_INT16": "32767",
		"TEST_INT32": "2147483647",
		"TEST_INT64": "9223372036854775807",
	})

	type args struct {
		TEST_INT   int
		TEST_INT8  int8
		TEST_INT16 int16
		TEST_INT32 int32
		TEST_INT64 int64
	}

	expect := &args{
		TEST_INT:   114514,
		TEST_INT8:  127,
		TEST_INT16: 32767,
		TEST_INT32: 2147483647,
		TEST_INT64: 9223372036854775807,
	}

	res := new(args)

	err := Parse(res)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expect, res)
}

func TestUintX(t *testing.T) {
	parseEnv(map[string]string{
		"TEST_UINT":   "12210",
		"TEST_UINT8":  "255",
		"TEST_UINT16": "65535",
		"TEST_UINT32": "4294967295",
		"TEST_UINT64": "18446744073709551615",
	})

	type args struct {
		TEST_UINT   uint
		TEST_UINT8  uint8
		TEST_UINT16 uint16
		TEST_UINT32 uint32
		TEST_UINT64 uint64
	}

	expect := &args{
		TEST_UINT:   12210,
		TEST_UINT8:  255,
		TEST_UINT16: 65535,
		TEST_UINT32: 4294967295,
		TEST_UINT64: 18446744073709551615,
	}

	res := new(args)

	err := Parse(res)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expect, res)
}

func TestFloatX(t *testing.T) {
	parseEnv(map[string]string{
		"TEST_FLOAT32": "3.1415926",
		"TEST_FLOAT64": "3.1415926535",
	})

	type args struct {
		TEST_FLOAT32 float32
		TEST_FLOAT64 float64
	}

	expect := &args{
		TEST_FLOAT32: 3.1415926,
		TEST_FLOAT64: 3.1415926535,
	}

	res := new(args)

	err := Parse(res)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expect, res)
}

func TestString(t *testing.T) {
	parseEnv(map[string]string{
		"TEST_STRING": "test",
	})

	type args struct {
		TEST_STRING string
	}

	expect := &args{
		TEST_STRING: "test",
	}

	res := new(args)

	err := Parse(res)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expect, res)
}

func TestStruct(t *testing.T) {
	parseEnv(map[string]string{
		"TEST_STRUCT_STRING": "test",
		"TEST_STRUCT_INT":    "114514",
	})

	type inline struct {
		TEST_STRUCT_STRING string
		TEST_STRUCT_INT    int
	}

	type args struct {
		Inline inline
	}

	expect := &args{
		Inline: inline{
			TEST_STRUCT_STRING: "test",
			TEST_STRUCT_INT:    114514,
		},
	}

	res := new(args)

	err := Parse(res)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expect, res)
}

func TestArrayAndSlice(t *testing.T) {
	parseEnv(map[string]string{
		"TEST_ARRAY_1": "1,3,5",
		"TEST_ARRAY_2": "2.0,4.0,6.0",
		"TEST_ARRAY_3": "true,false,true",
		"TEST_SLICE_1": "1,3,5",
		"TEST_SLICE_2": "2.0,4.0,6.0",
		"TEST_SLICE_3": "true,false,true",
	})

	type args struct {
		TEST_ARRAY_1 [3]int
		TEST_ARRAY_2 [3]float64
		TEST_ARRAY_3 [3]bool
		TEST_SLICE_1 []int
		TEST_SLICE_2 []float64
		TEST_SLICE_3 []bool
	}

	expect := &args{
		TEST_ARRAY_1: [3]int{1, 3, 5},
		TEST_ARRAY_2: [3]float64{2.0, 4.0, 6.0},
		TEST_ARRAY_3: [3]bool{true, false, true},
		TEST_SLICE_1: []int{1, 3, 5},
		TEST_SLICE_2: []float64{2.0, 4.0, 6.0},
		TEST_SLICE_3: []bool{true, false, true},
	}

	res := new(args)

	err := Parse(res)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expect, res)
}

func TestInterface(t *testing.T) {
	parseEnv(map[string]string{
		"TEST_INTERFACE": "114514,1919810",
	})

	type args struct {
		TEST_INTERFACE interface{}
	}

	expect := &args{TEST_INTERFACE: "114514,1919810"}

	res := new(args)

	err := Parse(res)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expect, res)
}

func TestZeroField(t *testing.T) {
	parseEnv(map[string]string{
		"TEST_ZERO_FIELD_ARRAY":     "114514,1919810",
		"TEST_ZERO_FIELD_SLICE":     "114514,1919810",
		"TEST_ZERO_FIELD_INTERFACE": "114514,1919810",
		"TEST_ZERO_FIELD_INT":       "12210",
	})

	type Inline struct {
		TEST_ZERO_FIELD_ARRAY     [2]int
		TEST_ZERO_FIELD_SLICE     []int
		TEST_ZERO_FIELD_INTERFACE interface{}
	}

	type args struct {
		TEST_ZERO_FIELD_INT int
		Inline
	}

	expect := &args{
		TEST_ZERO_FIELD_INT: 12210,
		Inline: Inline{
			TEST_ZERO_FIELD_ARRAY:     [2]int{114514, 1919810},
			TEST_ZERO_FIELD_SLICE:     []int{114514, 1919810},
			TEST_ZERO_FIELD_INTERFACE: "114514,1919810",
		},
	}

	res := &args{
		TEST_ZERO_FIELD_INT: 123456,
		Inline: Inline{
			TEST_ZERO_FIELD_ARRAY:     [2]int{123, 456},
			TEST_ZERO_FIELD_SLICE:     []int{123, 456},
			TEST_ZERO_FIELD_INTERFACE: "🇺🇦",
		},
	}

	parser, err := NewParser(&Config{
		ZeroFields: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	err = parser.Parse(res)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expect, res)
}

func TestTagName(t *testing.T) {
	parseEnv(map[string]string{
		"TEST_TAG_NAME_STRING": "Hexa",
		"TEST_TAG_NAME_INT":    "333333",
	})

	type args struct {
		Creator string `emp:"TEST_TAG_NAME_STRING"`
		Gulu    int    `emp:"name:TEST_TAG_NAME_INT"`
	}

	expect := &args{
		Creator: "Hexa",
		Gulu:    333333,
	}

	res := new(args)

	err := Parse(res)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expect, res)
}

func TestPrefix(t *testing.T) {
	parseEnv(map[string]string{
		"LOVELY_CREATOR":     "Hexa",
		"KALA_GULU":          "333333",
		"DING_DONG_DUANG":    "biubiubiu",
		"DING_DONG_BANGBANG": "gulugulu",
	})

	type inline struct {
		Duang    string `emp:"DUANG"`
		BANGBANG string
	}

	type args struct {
		CREATOR string `emp:"prefix:LOVELY_"`
		GULU    int    `emp:"prefix:KALA_"`
		Inline  inline `emp:"prefix:DING_DONG_"`
	}

	expect := &args{
		CREATOR: "Hexa",
		GULU:    333333,
		Inline: inline{
			Duang:    "biubiubiu",
			BANGBANG: "gulugulu",
		},
	}

	res := new(args)

	err := Parse(res)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expect, res)
}

func TestAllowEmpty(t *testing.T) {
	parseEnv(map[string]string{
		"TEST_ALLOW_EMPTY_INT": "114514",
	})

	type args struct {
		TEST_ALLOW_EMPTY_STRING string
		TEST_ALLOW_EMPTY_INT    int
	}

	expect := &args{
		TEST_ALLOW_EMPTY_STRING: "",
		TEST_ALLOW_EMPTY_INT:    114514,
	}

	res := new(args)

	parser, err := NewParser(&Config{
		AllowEmpty: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	err = parser.Parse(res)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expect, res)
}

func TestDefault(t *testing.T) {
	type args struct {
		TEST_DEFUALT_STRING string `emp:"default:LOVELY_CUTE_HEXAGRAM"`
		TEST_DEFUALT_INT    int    `emp:"default:333333"`
	}

	expect := &args{
		TEST_DEFUALT_STRING: "LOVELY_CUTE_HEXAGRAM",
		TEST_DEFUALT_INT:    333333,
	}

	res := new(args)

	err := Parse(res)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expect, res)
}

func TestDirectDefault(t *testing.T) {
	parseEnv(map[string]string{
		"TEST_DIRECT_DEFAULT_INT":    "131313",
		"TEST_DIRECT_DEFAULT_STRING": "MWE_MIAO_NYA",
	})

	type args struct {
		TEST_DIRECT_DEFAULT_STRING string `emp:"default:LOVELY_CUTE_HEXAGRAM"`
		TEST_DIRECT_DEFAULT_INT    int    `emp:"default:232323"`
	}

	expect := &args{
		TEST_DIRECT_DEFAULT_STRING: "LOVELY_CUTE_HEXAGRAM",
		TEST_DIRECT_DEFAULT_INT:    232323,
	}

	res := new(args)

	parser, err := NewParser(&Config{
		DirectDefault: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	err = parser.Parse(res)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expect, res)
}

func TestParseStringToArrayAndSlice(t *testing.T) {
	parseEnv(map[string]string{
		"TEST_ARRAY_1": "",
		"TEST_ARRAY_2": "[2.0,4.0,6.0]",
		"TEST_ARRAY_3": "[true,false,true]",
		"TEST_SLICE_1": "[1,3,5]",
		"TEST_SLICE_2": "[2.0,4.0,6.0]",
		"TEST_SLICE_3": "[true,false,true]",
	})

	type args struct {
		TEST_ARRAY_1 [3]int
		TEST_ARRAY_2 [3]float64
		TEST_ARRAY_3 [3]bool
		TEST_SLICE_1 []int
		TEST_SLICE_2 []float64
		TEST_SLICE_3 []bool
	}

	expect := &args{
		TEST_ARRAY_1: [3]int{},
		TEST_ARRAY_2: [3]float64{2.0, 4.0, 6.0},
		TEST_ARRAY_3: [3]bool{true, false, true},
		TEST_SLICE_1: []int{1, 3, 5},
		TEST_SLICE_2: []float64{2.0, 4.0, 6.0},
		TEST_SLICE_3: []bool{true, false, true},
	}

	res := new(args)

	parser, err := NewParser(&Config{
		AllowEmpty: true,
		ParseStringToArrayAndSlice: func(s string) []string {
			if s == "" {
				return []string{}
			}
			s = strings.TrimPrefix(s, "[")
			s = strings.TrimSuffix(s, "]")
			return strings.Split(s, ",")
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	err = parser.Parse(res)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expect, res)
}

func TestIgnorePrivate(t *testing.T) {
	parseEnv(map[string]string{
		"TEST_IGNORE_LOVELY_CREATOR":     "Hexa",
		"TEST_IGNORE_KALA_GULU":          "333333",
		"TEST_IGNORE_DING_DONG_DUANG":    "biubiubiu",
		"TEST_IGNORE_DING_DONG_BANGBANG": "gulugulu",
	})

	type inline struct {
		duang    string
		BANGBANG string
	}

	type args struct {
		LOVELY_CREATOR string `emp:"prefix:TEST_IGNORE_"`
		KALA_GULU      int    `emp:"prefix:TEST_IGNORE_"`
		Inline         inline `emp:"prefix:TEST_IGNORE_DING_DONG_"`
		inline         inline `emp:"prefix:TEST_IGNORE_DING_DONG_"`
	}

	expect := &args{
		LOVELY_CREATOR: "Hexa",
		KALA_GULU:      333333,
		Inline: inline{
			duang:    "",
			BANGBANG: "gulugulu",
		},
		inline: inline{
			duang:    "",
			BANGBANG: "",
		},
	}

	res := new(args)

	err := Parse(res)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expect, res)
}

func TestIgnore(t *testing.T) {
	parseEnv(map[string]string{
		"TEST_STRING": "test",
	})

	type args struct {
		TEST_STRING string `emp:"-"`
	}

	expect := &args{
		TEST_STRING: "",
	}

	res := new(args)

	err := Parse(res)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expect, res)
}

func TestAutoPrefix(t *testing.T) {
	parseEnv(map[string]string{
		"TEST_AUTO_PREFIX_A114514": "1919810",
	})

	type inline struct {
		A114514 int
	}

	type args struct {
		TEST_AUTO_PREFIX_ inline
	}

	expect := &args{
		inline{
			1919810,
		},
	}

	res := new(args)

	parser, err := NewParser(&Config{
		AutoPrefix: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	err = parser.Parse(res)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expect, res)
}

func TestMarshal(t *testing.T) {
	type Inline struct {
		STRING string
		ARRAY  []string
	}

	type args struct {
		BOOL bool
		INT  int
		Inline
	}

	expect := `
BOOL=false
INT=10000
STRING=gogogo
ARRAY=1,3,5
`

	input := &args{
		BOOL: false,
		INT:  10000,
		Inline: Inline{
			STRING: "gogogo",
			ARRAY:  []string{"1", "3", "5"},
		},
	}

	res, err := Marshal(input)
	if err != nil {
		panic(err)
	}

	fmt.Println(res)
	assert.Equal(t, expect, res)
}
