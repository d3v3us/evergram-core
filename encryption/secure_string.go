package encryption

import (
	"math/rand"
	"time"
)

// ISecureString defines the interface for a secure string.
type ISecureString interface {
	SetKey(key int)
	Set(value []byte) ISecureString
	Get() []byte
	IsEqual(other ISecureString) bool
}

// SecureString represents a securely stored string.
type SecureString struct {
	key         int
	realValue   []byte
	fakeValue   []byte
	initialized bool
}

// NewSecureString creates a new SecureString.
func NewSecureString(value string) ISecureString {
	s := &SecureString{
		key:         DefaultKey,
		realValue:   []byte(value),
		fakeValue:   []byte(value),
		initialized: false,
	}

	s.initialize()

	return s
}

// Initialize initializes the SecureString.
func (s *SecureString) initialize() {
	if !s.initialized {
		s.realValue = s.xor(s.realValue, s.key)
		s.initialized = true
	}
}

// SetKey sets the encryption key.
func (s *SecureString) SetKey(key int) {
	s.key = key
}

// RandomizeKey generates a random encryption key.
func (s *SecureString) randomizeKey() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	val, err := s.decrypt()
	if err != nil {
		panic(err)
	}
	s.realValue = val
	s.key = rand.Intn(int(^uint(0) >> 1))
	s.realValue = s.xor(s.realValue, s.key)
}

// XOR performs an XOR operation on the value with the key.
func (s *SecureString) xor(value []byte, key int) []byte {
	res := make([]byte, len(value))

	for i, v := range value {
		res[i] = v ^ byte(key)
	}

	return res
}

// Get returns the decrypted string.
func (s *SecureString) Get() []byte {
	val, err := s.decrypt()
	if err != nil {
		return nil
	}
	return val
}

// Set sets the secure string value.
func (s *SecureString) Set(value []byte) ISecureString {
	s.realValue = s.xor(value, s.key)
	return s
}

// Decrypt decrypts the secure string.
func (s *SecureString) decrypt() ([]byte, error) {
	if !s.initialized {
		s.key = DefaultKey
		s.fakeValue = nil
		s.realValue = s.xor(nil, 0)
		s.initialized = true

		return nil, nil
	}

	res := s.xor(s.realValue, s.key)

	return res, nil
}

// IsEqual checks if two secure strings are equal.
func (s *SecureString) IsEqual(other ISecureString) bool {
	if s.key != other.(*SecureString).key {
		return string(s.xor(s.realValue, s.key)) == string(s.xor(other.(*SecureString).realValue, other.(*SecureString).key))
	}

	return string(s.realValue) == string(other.(*SecureString).realValue)
}

// Constants
const (
	DefaultKey = 12345
)
