// Maybe Monad represents the presence `Just` or an absence `Nothing` of a
// value. It is commonly used to handle computations that may or may not return
// a result, or computations that can potentially fail.
//
// Inspiration and code taken from:
// * https://github.com/erikjuhani/go-fp/blob/main/maybe/maybe.go
// * https://raw.githubusercontent.com/pmorelli92/maybe/main/maybe.go
package mo

import (
	"encoding/json"
	"reflect"
)

// Maybe monad data type representation. May or may not contain a pointer
// value. Nothing is represented as a `nil` value internally.
type Maybe[T any] struct {
	value T
	valid bool
}

func (m Maybe[T]) HasValue() bool {
	return m.valid
}

// Value returns the value of the Maybe. It does not protect against nil, so how can we do that?
func (m Maybe[T]) Value() T {
	return m.value
}

// ValueOr uses the value if set in the Maybe, otherwise uses the passed value
func (m Maybe[T]) ValueOr(v T) T {
	if m.valid {
		return m.value
	}
	return v
}

func (m *Maybe[T]) UnmarshalJSON(data []byte) error {
	var t *T
	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}

	if t != nil {
		*m = Just(*t)
	}

	return nil
}

func (m Maybe[T]) MarshalJSON() ([]byte, error) {
	var t *T

	if m.valid {
		t = &m.value
	}

	return json.Marshal(t)
}

// Just is the return operation for Maybe monad that returns the representation
// of existence of a value.
func Just[T any](v T) Maybe[T] {
	if reflect.ValueOf(v).Kind() == reflect.Ptr {
		rv := reflect.ValueOf(v)
		if rv.IsNil() {
			return Nothing[T]()
		}
		x := rv.Elem().Interface().(T)
		return Maybe[T]{value: x, valid: true}
	}
	return Maybe[T]{value: v, valid: true}
}

// Nothing is the return operation for Maybe monad that returns the representation
// of absence of a value.
func Nothing[T any]() Maybe[T] {
	return Maybe[T]{valid: false}
}

// From is the return operation for Maybe monad that returns either Just a or
// Nothing. Intended to be used with Go functions that return tuple as `val, ok`.
func From[T any](val T, ok ...bool) Maybe[T] {
	// TOOD: understand this "Ok" logic
	if len(ok) > 0 && !ok[0] {
		return Nothing[T]()
	}
	if reflect.ValueOf(val).Kind() == reflect.Ptr {
		rv := reflect.ValueOf(val)
		if rv.IsNil() {
			return Nothing[T]()
		}
		x := rv.Elem().Interface().(T)
		return Maybe[T]{value: x, valid: true}
	}
	return Just(val)
}

// Map or the "bind" function takes the contents of the Maybe monad and passes
// it to function `f` as a parameter. The function `f` returns a new Maybe
// monad as a result.
func Map[A, B any](f func(A) B) func(Maybe[A]) Maybe[B] {
	return func(m Maybe[A]) Maybe[B] {
		if m.valid {
			return Just(f(m.value))
		}
		return Nothing[B]()
	}
}

// Fmap or also known as `bind` function lets non-monadic function `f` to
// operate on the contents of monad m a, and lifts the value to a new domain
// (Maybe a -> Maybe b).
func Fmap[A, B any](f func(A) Maybe[B]) func(Maybe[A]) Maybe[B] {
	return func(m Maybe[A]) Maybe[B] {
		if m.valid {
			return f(m.value)
		}
		return Nothing[B]()
	}
}

// Match matches Maybe monad depending of it's current state and returns the
// value determined by the return type of b.
func Match[A, B any](Nothing func() B, Just func(A) B) func(Maybe[A]) B {
	return func(m Maybe[A]) B {
		if m.valid {
			return Just(m.value)
		}
		return Nothing()
	}
}
