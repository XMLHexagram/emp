package main

import (
	"fmt"
	"github.com/XMLHexagram/emp"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	//res, err := emp.Marshal(&Config{})
	//if err != nil {
	//	panic(err)
	//}
	//
	//fmt.Println(res)

	env, err := ParserEnv()
	if err != nil {
		panic(err)
	}

	// every environment you need will be filled into `env` struct
	fmt.Printf("%#v\n", env)
}

func ParserEnv() (*Config, error) {
	config := new(Config)

	err := emp.Parse(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
