package optional

// Optional represents a value that may or may not be present
type Optional[T any] struct {
	value   T
	present bool
}

// Some creates an Optional with a value
func Some[T any](value T) Optional[T] {
	return Optional[T]{
		value:   value,
		present: true,
	}
}

// None creates an empty Optional
func None[T any]() Optional[T] {
	return Optional[T]{
		present: false,
	}
}

// IsPresent returns true if the Optional contains a value
func (o Optional[T]) IsPresent() bool {
	return o.present
}

// IsEmpty returns true if the Optional is empty
func (o Optional[T]) IsEmpty() bool {
	return !o.present
}

// Get returns the value and a boolean indicating if it was present
func (o Optional[T]) Get() (T, bool) {
	return o.value, o.present
}

// Unwrap returns the value, panicking if it's not present
func (o Optional[T]) Unwrap() T {
	if !o.present {
		panic("called Unwrap on an empty Optional")
	}
	return o.value
}

// UnwrapOr returns the value if present, otherwise returns the default value
func (o Optional[T]) UnwrapOr(defaultValue T) T {
	if o.present {
		return o.value
	}
	return defaultValue
}

// Map applies a function to the contained value if present
func Map[T, U any](o Optional[T], f func(T) U) Optional[U] {
	if o.present {
		return Some(f(o.value))
	}
	return None[U]()
}

// FlatMap applies a function that returns an Optional to the contained value
func FlatMap[T, U any](o Optional[T], f func(T) Optional[U]) Optional[U] {
	if o.present {
		return f(o.value)
	}
	return None[U]()
}
