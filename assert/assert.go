package assert

import (
	"strings"
)

func Assert(cond bool, msg string) {
	if !cond {
		panic(msg)
	}
}

func AssertNotNil(a any, msg ...string) {
	if a == nil {
		panic(strings.Join(append([]string{"value is nil"}, msg...), ", "))
	}
}

func AssertNotEmpty(s string, msg string) {
	if s == "" {
		panic("zero value: " + msg)
	}
}
