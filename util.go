package emp

import (
	"github.com/XMLHexagram/emp/empErr"
	"os"
	"reflect"
	"strings"
)

func parseTagString(tagString string) (prefix string, name string, default_ string, isIgnore bool) {
	tagParts := strings.Split(tagString, ",")
	for _, tagPart := range tagParts {
		if tagPart == "-" {
			isIgnore = true
		} else if strings.HasPrefix(tagPart, "prefix=") {
			prefix = strings.TrimPrefix(tagPart, "prefix=")
		} else if strings.HasPrefix(tagPart, "name=") {
			name = strings.TrimPrefix(tagPart, "name=")
		} else if strings.HasPrefix(tagPart, "default=") {
			default_ = strings.TrimPrefix(tagPart, "default=")
		}
	}
	return prefix, name, default_, isIgnore
}

func getKind(val reflect.Value) reflect.Kind {
	kind := val.Kind()

	switch {
	case kind >= reflect.Uint && kind <= reflect.Uint64:
		return reflect.Uint
	default:
		return kind
	}
}

func ParseStringToSlice(s string) []string {
	if s == "" {
		return []string{}
	}
	return strings.Split(s, ",")
}

func getEnvString(key string, default_ string, directDefault bool, allowEmpty bool) (envString string, err error) {
	if directDefault {
		envString = default_
	} else {
		envString = os.Getenv(key)
		if envString == "" {
			envString = default_
		}
		if envString == "" && !allowEmpty {
			return "", empErr.NotAllowEmptyEnvError.New().Wrap("miss environment key: " + key)
		}
	}
	return envString, nil
}
