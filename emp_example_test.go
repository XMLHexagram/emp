package emp

import (
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
