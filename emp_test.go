package emp

import (
	"github.com/joho/godotenv"
	"path/filepath"
	"testing"
)

func Test1(t *testing.T) {
	godotenv.Parse()
	err := godotenv.Load(filepath.Join("testData", ".env"))
	if err != nil {
		t.Fatal(err)
	}

	testStruct := &struct {
		HELLO []interface{}
	}{}

	err = Parse(testStruct)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(testStruct)
}
