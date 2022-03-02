package main

import (
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParserEnv(t *testing.T) {
	expect := &Config{
		Server: Server{
			Http: Http{
				Port:    ":12210",
				Timeout: 200,
			},
		},
		Db: Db{
			Driver: "postgres",
			DSN:    "postgres://postgres:postgres@localhost:5432/postgres",
		},
		Cache: Cache{
			Driver: "redis",
			DSN:    "redis://localhost:6379/0",
		},
		Log: Log{
			Level:       "info",
			ToStd:       true,
			LogRotate:   true,
			Development: true,
			Sampling:    true,
			Rotate: Rotate{
				Filename:   "emp.log",
				MaxSize:    200,
				MaxAge:     10,
				MaxBackups: 10},
		},
	}

	err := godotenv.Load()
	if err != nil {
		t.Fatal(err)
	}

	res, err := ParserEnv()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expect, res)
}
