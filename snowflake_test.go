package snowflake

import (
	"github.com/bmizerany/assert"
	"testing"
)

func TestSnowflake_Base64(t *testing.T) {
	t.Run("max_value", func(t *testing.T) {
		s := Snowflake(maxValueBits(63))
		assert.Equal(t, "7//////////", s.Base64())
	})
	t.Run("min_value", func(t *testing.T) {
		s := Snowflake(0)
		assert.Equal(t, "0", s.Base64())
	})
}

func TestParseBase64(t *testing.T) {
	t.Run("max_value", func(t *testing.T) {
		s, err := ParseBase64("7//////////")
		assert.Equal(t, nil, err)
		assert.Equal(t, Snowflake(maxValueBits(63)), s)
	})
	t.Run("min_value", func(t *testing.T) {
		s, err := ParseBase64("0")
		assert.Equal(t, nil, err)
		assert.Equal(t, Snowflake(0), s)
	})
	t.Run("middle_value", func(t *testing.T) {
		s, err := ParseBase64("3//////////")
		assert.Equal(t, nil, err)
		assert.Equal(t, Snowflake(4611686018427387903), s)
	})
	t.Run("half_set", func(t *testing.T) {
		s, err := ParseBase64("400000") // max uint32 (half of the bits are set)
		assert.Equal(t, nil, err)
		assert.Equal(t, Snowflake(4294967296), s)
	})
	t.Run("invalid", func(t *testing.T) {
		_, err := ParseString("4294967296a.")
		assert.NotEqual(t, nil, err)
	})
}

func TestSnowflake_String(t *testing.T) {
	t.Run("max_value", func(t *testing.T) {
		s := Snowflake(maxValueBits(63))
		assert.Equal(t, "9223372036854775807", s.String())
	})
	t.Run("min_value", func(t *testing.T) {
		s := Snowflake(0)
		assert.Equal(t, "0", s.String())
	})
}

func TestParseString(t *testing.T) {
	t.Run("max_value", func(t *testing.T) {
		s, err := ParseString("9223372036854775807")
		assert.Equal(t, err, nil)
		assert.Equal(t, Snowflake(9223372036854775807), s)
	})
	t.Run("min_value", func(t *testing.T) {
		s, err := ParseString("0")
		assert.Equal(t, nil, err)
		assert.Equal(t, Snowflake(0), s)
	})
	t.Run("middle_value", func(t *testing.T) {
		s, err := ParseString("4611686018427387903")
		assert.Equal(t, nil, err)
		assert.Equal(t, Snowflake(4611686018427387903), s)
	})
	t.Run("half_set", func(t *testing.T) {
		s, err := ParseString("4294967296") // max uint32 (half of the bits are set)
		assert.Equal(t, nil, err)
		assert.Equal(t, Snowflake(4294967296), s)
	})
	t.Run("invalid", func(t *testing.T) {
		_, err := ParseString("4294967296a")
		assert.NotEqual(t, nil, err)
	})
}

func TestSnowflake_MarshalJSON(t *testing.T) {
	t.Run("max_value", func(t *testing.T) {
		s := Snowflake(maxValueBits(63))
		b, err := s.MarshalJSON()
		assert.Equal(t, nil, err)
		assert.Equal(t, "\"9223372036854775807\"", string(b))
	})
	t.Run("min_value", func(t *testing.T) {
		s := Snowflake(0)
		b, err := s.MarshalJSON()
		assert.Equal(t, nil, err)
		assert.Equal(t, "\"0\"", string(b))
	})
}

func TestSnowflake_UnmarshalJSON(t *testing.T) {
	t.Run("max_value", func(t *testing.T) {
		var s Snowflake
		err := s.UnmarshalJSON([]byte("\"9223372036854775807\""))
		assert.Equal(t, err, nil)
		assert.Equal(t, Snowflake(9223372036854775807), s)
	})
	t.Run("min_value", func(t *testing.T) {
		var s Snowflake
		err := s.UnmarshalJSON([]byte("\"0\""))
		assert.Equal(t, nil, err)
		assert.Equal(t, Snowflake(0), s)
	})
	t.Run("middle_value", func(t *testing.T) {
		var s Snowflake
		err := s.UnmarshalJSON([]byte("\"4611686018427387903\""))
		assert.Equal(t, nil, err)
		assert.Equal(t, Snowflake(4611686018427387903), s)
	})
	t.Run("half_set", func(t *testing.T) {
		var s Snowflake
		err := s.UnmarshalJSON([]byte("\"4294967296\"")) // max uint32 (half of the bits are set)
		assert.Equal(t, nil, err)
		assert.Equal(t, Snowflake(4294967296), s)
	})
	t.Run("invalid", func(t *testing.T) {
		var s Snowflake
		err := s.UnmarshalJSON([]byte("\"4294967296a\""))
		assert.NotEqual(t, nil, err)
	})
}

func BenchmarkParseBase64(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = ParseBase64("7//////////")
	}
}
