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
