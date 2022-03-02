package emp

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReadmeExample(t *testing.T) {
	parseEnv(map[string]string{
		"SOME_ENV":   "1221",
		"SOME_ENV_1": "hello",
		"SOME_ENV_2": "lovely,cute,Hexagram",
	})

	type EnvModel struct {
		SOME_ENV   int
		SOME_ENV_1 string
		SOME_ENV_2 []string
	}

	expect := &EnvModel{
		SOME_ENV:   1221,
		SOME_ENV_1: "hello",
		SOME_ENV_2: []string{"lovely", "cute", "Hexagram"},
	}

	res := new(EnvModel)

	err := Parse(res)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, expect, res)
}

func TestReadmeExample2(t *testing.T) {
	type EnvModel struct {
		JWT_SECERT string
		JWT_EXPIRE int
		REDIS_URL  string
		SERVER_    struct {
			PORT         string
			HTTP_TIMEOUT int
		}
		// whatever type you want, but not map
	}

	res, err := Marshal(&EnvModel{})
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(res)

	parser, err := NewParser(&Config{
		AutoPrefix: true,
	})
	if err != nil {
		t.Fatal(err)
	}

	res, err = parser.Marshal(&EnvModel{})
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(res)

	type EnvModel1 struct {
		JwtSecret string `emp:"JWT_SECRET"`
		JwtExpire int    `emp:"JWT_EXPIRE"`
		RedisUrl  string `emp:"REDIS_URL"`
		Server    struct {
			Port        string `emp:"SERVER_PORT"`
			HttpTimeout int    `emp:"SERVER_HTTP_TIMEOUT"`
		} `emp:"prefix:SERVER_"`
		// whatever type you want, but not map
	}

	res, err = Marshal(&EnvModel1{})
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(res)
}
