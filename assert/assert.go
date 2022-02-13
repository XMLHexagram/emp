package assert

import (
	"reflect"
)

func NotNil() {

}

func Equal(actual, expected interface{}) {
	if reflect.DeepEqual(actual, expected) {
		panic("assertion failed")
	}
}
