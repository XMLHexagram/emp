# emp 

[![Go Report Card](https://goreportcard.com/badge/github.com/XMLHexagram/emp)](https://goreportcard.com/report/github.com/XMLHexagram/emp)
[![Godoc](https://godoc.org/github.com/XMLHexagram/emp?status.svg)](https://pkg.go.dev/github.com/XMLHexagram/emp)
[![Test](https://github.com/XMLHexagram/emp/actions/workflows/go.yml/badge.svg)](https://github.com/XMLHexagram/emp/actions/workflows/go.yml)
![Go Version](https://img.shields.io/badge/go%20version-%3E=1.14-61CFDD.svg)

## What is emp?

emp is a Go library for parsing environment variable to structures, while providing helpful error handling.

## Why emp?

Directly read environment variables in Go will case a lot of repeat code like this:

```go
package main

import "os"
    
func main() {
    envStringValue := os.Getenv("SOME_ENV")
	
    if envStringValue == "" {
        panic("missing SOME_ENV") 
        // OR: envStringValue = "default"
    }
    
    // convert string to the type you want
    envValue := sthConvert(envValue)
}    
```

It's a bit complicated and unclear, isn't it?

Let's try emp!

```go
package main

import "github.com/XMLHexagram/emp"

func main() {
    type EnvModel struct {
        SOME_ENV   int
        SOME_ENV_1 string
        SOME_ENV_2 []string
        // whatever type you want, but not map
    }

    envModel := new(EnvModel)
    
    err := Parse(envModel)
    
    if err != nil {
        panic(err)
    }
}
```

Now, `SOME_ENV` is parsed to `envModel`.

It is simpler and easier to understand than the previous one, right?

What, no? Then you need to see [Usage](#Usage) below.

## Installation

Standard `go get`:

```
$ go get github.com/XMLHexagram/emp
```

## Usage

For full usage and examples see the [Godoc](http://godoc.org/github.com/XMLHexagram/emp).

Or you can see [Real Project Example](https://github.com/XMLHexagram/emp/tree/main/example),
[Test Case](https://github.com/XMLHexagram/emp/blob/main/emp_test.go) and
[Example Test Case](https://github.com/XMLHexagram/emp/blob/main/emp_example_test.go)

Here is a short introduction of emp:

### Parse Environment variables to struct

The easiest way to use emp is to use `Parse` function.

First, define an empty struct:

```go
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
```

Then, Let's use `emp.Marshal` to see the environment variables emp will looking for.

```go
res, err := emp.Marshal(&EnvModel{})

fmt.Println(res)
/*
JWT_SECERT=
JWT_EXPIRE=0
REDIS_URL=
// you want SERVER_PORT and SERVER_HTTP_TIMEOUT, right?
PORT=
HTTP_TIMEOUT=0
*/
```

Good, but not good enough. Maybe you want `SERVER_PORT` and `SERVER_HTTP_TIMEOUT`.

Let's try `AutoPrefix`:

```go
parser, _ := emp.NewParser(&Config{
    AutoPrefix: true,
})

res, _ = parser.Marshal(&EnvModel{})

fmt.Println(res)
/*
JWT_SECERT=
JWT_EXPIRE=0
REDIS_URL=
// the prefix is `SERVER_` now!
SERVER_PORT=
SERVER_HTTP_TIMEOUT=0
*/
```

Much better, but can emp do something more?

Of course can, let's try `field tag` to customize the `prefix` and `environment variable` emp looking for:

```go
type EnvModel struct {
    JwtSecret string `emp:"JWT_SECRET"`
    JwtExpire int    `emp:"JWT_EXPIRE"`
    RedisUrl  string `emp:"REDIS_URL"`
    Server    struct {
        Port        string `emp:"SERVER_PORT"`
        HttpTimeout int    `emp:"SERVER_HTTP_TIMEOUT"`
    } `emp:"prefix:SERVER_"`
    // whatever type you want, but not map
}

res, _ := emp.Marshal(&EnvModel1{})

fmt.Println(res)
/*
JWT_SECRET=
JWT_EXPIRE=0
REDIS_URL=
SERVER_SERVER_PORT=
SERVER_SERVER_HTTP_TIMEOUT=0
*/
```

COOL! This struct defines looks perfect. Now, we only need one more step.

Finally, call `Parse` function:

```go
envModel := new(EnvModel)

_ := emp.Parse(envModel)
```

envModel is now filled with environment variables you need. 🎉

That's All? 

No, emp also provides more features to customize your environment variables parsing.

See [emp doc](https://godoc.org/github.com/XMLHexagram/emp) for more details.

## Q & A

### Why is it called "emp"?

A: emp => **E**nviron**M**ent variable **P**arser

### Does emp support case sensitive keys?

Yes, but maybe not in the future.

As we all know, environment variables in most time are ALL CAPS. So if someone need this feature, I will add it as a 
configurable option.

### How to set environment variables?

There is an easy way to set environment variables by use [godotenv](https://github.com/joho/godotenv)

### Support for map type?

No, the only not supported native type of Go is map. Sincerely, I have no idea on how to support it.
If you have any idea, please [open an issue](https://github.com/XMLHexagram/emp/issues/new).

## Thanks

Thanks to [mitchellh/mapstructure](https://github.com/mitchellh/mapstructure) and [spf13/viper](https://github.
com/spf13/viper), Your repos have inspired me, and I used some code from your repo.

## More Feature or Bug Report

Feel free to [open an issue](https://github.com/XMLHexagram/emp/issues/new) when you want to report a bug, feature 
or ask for help. ✨