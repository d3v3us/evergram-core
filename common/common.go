package common

import (
	"strconv"
	"strings"
)

type Extendable[T any] struct {
	Value T
}

func ExtOf[T any](value T) *Extendable[T] {
	return &Extendable[T]{Value: value}
}

type Else[T any] interface {
	ElseDo(fn func() T) T
	Else(value T) T
}

type Then[T any] interface {
	ThenDo(fn func() T) Else[T]
	Then(value T) Else[T]
}

type Condition[T any] struct {
	condition bool
	thenValue T
	thenFn    func() T
}

func When[T any](condition bool) Then[T] {
	return &Condition[T]{condition: condition}
}

func (c *Condition[T]) ThenDo(fn func() T) Else[T] {
	c.thenFn = fn
	return c
}

func (c *Condition[T]) Then(value T) Else[T] {
	c.thenValue = value
	return c
}

func (c *Condition[T]) ElseDo(fn func() T) T {
	if c.condition {
		return c.then()
	}

	return fn()
}

func (c *Condition[T]) Else(value T) T {
	if c.condition {
		return c.then()
	}

	return value
}

func (c *Condition[T]) then() T {
	if c.thenFn != nil {
		return c.thenFn()
	}
	return c.thenValue
}

// Location represents the location details associated with a habit.
type Location struct {
	Latitude  float64 // Latitude of the location.
	Longitude float64 // Longitude of the location.
	Address   string  // Address of the location.
	City      string  // City of the location.
	State     string  // State or province of the location.
	Country   string  // Country of the location.
}

// RecurrencePattern represents the pattern for a recurring habit.
type RecurrencePattern struct {
	Frequency     int    // Frequency of recurrence (e.g., every 2 days).
	FrequencyType string // Type of frequency (e.g., day, week, month, year).
}

func StrToInt(str string) int {
	// Delete any non-numeric character from the string
	str = strings.TrimFunc(str, func(r rune) bool {
		return r < '0' || r > '9'
	})

	// Convert the cleaned-up string to an integer
	n, _ := strconv.Atoi(str)
	return n
}

func StrToFloat(str string) float64 {
	// Delete any non-numeric character from the string
	str = strings.TrimFunc(str, func(r rune) bool {
		return r < '0' || r > '9'
	})
	// Convert the cleaned-up string to an integer
	n, _ := strconv.ParseFloat(str, 64)
	return n
}

func StrToByteArray(str string) []byte {
	bytes := make([]byte, len(str))
	copy(bytes, str)

	return bytes
}
