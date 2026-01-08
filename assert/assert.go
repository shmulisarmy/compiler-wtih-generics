package assert

import "fmt"

func Equal[T comparable](t T, expected T, message ...string) {
	if len(message) == 0 {
		message = append(message, fmt.Sprintf("expected %v, got %v", expected, t))
	}
	if t != expected {
		panic(message[0])
	}
}

func NotEqual[T comparable](t T, un_expected T) {
	if t == un_expected {
		panic(fmt.Sprintf("%v and %v are equal", t, un_expected))
	}
}

func Assert(t bool, message ...string) {
	if len(message) == 0 {
		message = append(message, "assert failed")
	}
	if !t {
		panic(message[0])
	}
}
