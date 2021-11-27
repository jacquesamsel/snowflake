package snowflake

import (
	"encoding/json"
	"errors"
	"reflect"
	"strconv"
	"strings"
)

var (
	baseChars        = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ+/")
	// ErrorInvalidChar is when a snowflake is when the string being parsed contains a character outside the allowed
	// range (i.e. a-zA-Z0-9+/ for Base64)
	ErrorInvalidChar = errors.New("invalid character")
)

// Snowflake is a 63-bit integer made up of 3 parts - a timestamp, Node ID and a sequence ID which is calculated by the
// Node responsible for issuing the Snowflake.
type Snowflake int64

// Base64 converts the snowflake into the mathematical 64th base (i.e. not the one defined in rfc4648). This allows for
// higher efficiency (in terms of result length) because the integer does not need to processed in groups.
// The maximum resulting string length should be 8 chars
func (s Snowflake) Base64() string {
	var sb strings.Builder
	sb.Grow(8) // bytes in 64 bits, should fit up to the maximum of a uint64 (////////).
	for s != 0 {
		sb.WriteRune(baseChars[s%64])
		s /= 64
	}
	if sb.Len() == 0 {
		return "0"
	}
	return reverse(sb.String())
}

// String converts the Snowflake to a base-10 string
func (s Snowflake) String() string {
	return strconv.FormatInt(int64(s), 10)
}

// MarshalJSON is the implementation of json.Marshaler
func (s Snowflake) MarshalJSON() ([]byte, error) {
	return []byte("\"" + s.String() + "\""), nil
}

// UnmarshalJSON is the implementation of json.Unmarshaler
func (s *Snowflake) UnmarshalJSON(b []byte) error {
	if b[0] != '"' || b[len(b)-1] != '"' || len(b) <= 2 {
		return &json.InvalidUnmarshalError{Type: reflect.TypeOf(s)}
	}
	i, err := strconv.ParseInt(string(b[1:len(b)-1]), 10, 64)
	*s = Snowflake(i)
	return err
}

//reverse reverses a string.
func reverse(s string) string {
	var sb strings.Builder
	str := []rune(s)
	for i := range str {
		sb.WriteRune(str[len(str)-1-i])
	}
	return sb.String()
}

// ParseBase64 parses a Snowflake from the 64th mathematical base.
func ParseBase64(s string) (Snowflake, error) {
	currentMultiplier := int64(1)
	var val int64
	for i := range s {
		var j int32
		c := int32(s[len(s)-1-i])
		// this chain of if statements is used, because it is faster than iterating through a []rune to find the index
		// of the character.
		if '0' <= c && c <= '9' {
			j = c - '0'
		} else if 'a' <= c && c <= 'z' {
			j = c - 'a' + 10
		} else if 'A' <= c && c <= 'Z' {
			j = c - 'A' + 36
		} else if c == '+' {
			j = 62
		} else if c == '/' {
			j = 63
		} else { // The character was not a valid 64th base character [1-9a-zA-Z+/]
			return 0, ErrorInvalidChar
		}
		val += currentMultiplier * int64(j)
		currentMultiplier *= 64
	}
	return Snowflake(val), nil
}

func ParseString(s string) (Snowflake, error) {
	i, err := strconv.ParseInt(s, 10, 64)
	return Snowflake(i), err
}
