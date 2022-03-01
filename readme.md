# emp 

[![Go Report Card](https://goreportcard.com/badge/github.com/XMLHexagram/emp)](https://goreportcard.com/report/github.com/XMLHexagram/emp)
[![Godoc](https://godoc.org/github.com/XMLHexagram/emp?status.svg)](https://pkg.go.dev/github.com/XMLHexagram/emp)

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

## Installation

Standard `go get`:

```
$ go get github.com/XMLHexagram/emp
```

## Usage & Example

For full usage and examples see the [Godoc](http://godoc.org/github.com/XMLHexagram/emp).

Or you can see [Example](https://github.com/XMLHexagram/emp/tree/main/example),
[Test Case](https://github.com/XMLHexagram/emp/blob/main/emp_test.go) and
[Example Test Case](https://github.com/XMLHexagram/emp/blob/main/emp_example_test.go)

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

Thanks to [mitchellh/mapstructure](https://github.com/mitchellh/mapstructure), I used code from your repo.

## More Feature or Bug Report

Feel free to [open an issue](https://github.com/XMLHexagram/emp/issues/new) when you want to report a bug, feature 
or ask for help. ✨