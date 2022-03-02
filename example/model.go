package main

type Config struct {
	Server `emp:"prefix:SERVER_"`
	Db     `emp:"prefix:DB_"`
	Cache  `emp:"prefix:CACHE_"`
	Log    `emp:"prefix:LOG_"`
}

type Http struct {
	Port    string `emp:"HTTP_PORT"`
	Timeout int    `emp:"HTTP_TIMEOUT"`
}

type Server struct {
	Http Http
}

type Db struct {
	Driver string `emp:"DRIVER"`
	DSN    string `emp:"DSN"`
}

//type DbMap map[string]Db

type Cache struct {
	Driver string `emp:"DRIVER"`
	DSN    string `emp:"DSN"`
}

//type CacheMap map[string]Cache

type Log struct {
	Level       string `emp:"LEVEL"`
	ToStd       bool   `emp:"TO_STD"`
	LogRotate   bool   `emp:"ROTATE"`
	Development bool   `emp:"DEVELOPMENT"`
	Sampling    bool   `emp:"SAMPLING"`
	Rotate      Rotate `emp:"prefix:ROTATE_"`
}

type Rotate struct {
	Filename   string `emp:"FILENAME"`
	MaxSize    int    `emp:"MAXSIZE"`
	MaxAge     int    `emp:"MAXAGE"`
	MaxBackups int    `emp:"MAXBACKUPS"`
}
