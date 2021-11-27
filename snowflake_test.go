package snowflake

import (
	"errors"
	"testing"
)

func TestSnowflake_Base64(t *testing.T) {
	t.Run("max_value", func(t *testing.T) {
		s := Snowflake(maxValueBits(63))
		if s.Base64() != "7//////////" {
			t.Fail()
		}
	})
	t.Run("min_value", func(t *testing.T) {
		s := Snowflake(0)
		if s.Base64() != "0" {
			t.Fail()
		}
	})
}

func TestParseBase64(t *testing.T) {
	t.Run("max_value", func(t *testing.T) {
		s, err := ParseBase64("7//////////")
		if err != nil {
			t.Error(err)
		}
		if s != Snowflake(maxValueBits(63)) {
			t.Fail()
		}
	})
	t.Run("min_value", func(t *testing.T) {
		s, err := ParseBase64("0")
		if err != nil {
			t.Error(err)
		}
		if s != Snowflake(0) {
			t.Fail()
		}
	})
	t.Run("middle_value", func(t *testing.T) {
		s, err := ParseBase64("3//////////")
		if err != nil {
			t.Error(err)
		}
		if s != Snowflake(4611686018427387903) {
			t.Fail()
		}
	})
	t.Run("half_set", func(t *testing.T) {
		s, err := ParseBase64("400000") // max uint32 (half of the bits are set)
		if err != nil {
			t.Error(err)
		}
		if s != Snowflake(4294967296) {
			t.Fail()
		}
	})
	t.Run("invalid", func(t *testing.T) {
		_, err := ParseString("4294967296a.")
		if !errors.Is(err, ParseError) {
			t.Fail()
		}
	})
}

func TestSnowflake_String(t *testing.T) {
	t.Run("max_value", func(t *testing.T) {
		s := Snowflake(maxValueBits(63))
		if s.String() != "9223372036854775807" {
			t.Fail()
		}
	})
	t.Run("min_value", func(t *testing.T) {
		s := Snowflake(0)
		if s.String() != "0" {
			t.Fail()
		}
	})
}

func TestParseString(t *testing.T) {
	t.Run("max_value", func(t *testing.T) {
		s, err := ParseString("9223372036854775807")
		if err != nil {
			t.Error(err)
		}
		if s != Snowflake(9223372036854775807) {
			t.Fail()
		}
	})
	t.Run("min_value", func(t *testing.T) {
		s, err := ParseString("0")
		if err != nil {
			t.Error(err)
		}
		if s != Snowflake(0) {
			t.Fail()
		}
	})
	t.Run("middle_value", func(t *testing.T) {
		s, err := ParseString("4611686018427387903")
		if err != nil {
			t.Error(err)
		}
		if s != Snowflake(4611686018427387903) {
			t.Fail()
		}
	})
	t.Run("half_set", func(t *testing.T) {
		s, err := ParseString("4294967296") // max uint32 (half of the bits are set)
		if err != nil {
			t.Error(err)
		}
		if s != Snowflake(4294967296) {
			t.Fail()
		}
	})
	t.Run("invalid", func(t *testing.T) {
		_, err := ParseString("4294967296a")
		if !errors.Is(err, ParseError) {
			t.Fail()
		}
	})
}

func TestSnowflake_MarshalJSON(t *testing.T) {
	t.Run("max_value", func(t *testing.T) {
		s := Snowflake(maxValueBits(63))
		b, err := s.MarshalJSON()
		if err != nil {
			t.Error(err)
		}
		if string(b) != "\"9223372036854775807\"" {
			t.Fail()
		}
	})
	t.Run("min_value", func(t *testing.T) {
		s := Snowflake(0)
		b, err := s.MarshalJSON()
		if err != nil {
			t.Error(err)
		}
		if string(b) != "\"0\"" {
			t.Fail()
		}
	})
}

func TestSnowflake_UnmarshalJSON(t *testing.T) {
	t.Run("max_value", func(t *testing.T) {
		var s Snowflake
		err := s.UnmarshalJSON([]byte("\"9223372036854775807\""))
		if err != nil {
			t.Fail()
		}
		if s != Snowflake(9223372036854775807) {
			t.Fail()
		}
	})
	t.Run("min_value", func(t *testing.T) {
		var s Snowflake
		err := s.UnmarshalJSON([]byte("\"0\""))
		if err != nil {
			t.Error(err)
		}
		if s != Snowflake(0) {
			t.Fail()
		}
	})
	t.Run("middle_value", func(t *testing.T) {
		var s Snowflake
		err := s.UnmarshalJSON([]byte("\"4611686018427387903\""))
		if err != nil {
			t.Error(err)
		}
		if s != Snowflake(4611686018427387903) {
			t.Fail()
		}
	})
	t.Run("half_set", func(t *testing.T) {
		var s Snowflake
		err := s.UnmarshalJSON([]byte("\"4294967296\"")) // max uint32 (half of the bits are set)
		if err != nil {
			t.Error(err)
		}
		if s != Snowflake(4294967296) {
			t.Fail()
		}
	})
	t.Run("invalid", func(t *testing.T) {
		t.Run("short", func(t *testing.T) {
			var s Snowflake
			err := s.UnmarshalJSON([]byte("a"))
			if !errors.Is(err, JSONUnmarshalError) {
				t.Fail()
			}
		})
		t.Run("char_outside_rng", func(t *testing.T) {
			var s Snowflake
			err := s.UnmarshalJSON([]byte("\"4294967296a\""))
			if !errors.Is(err, JSONUnmarshalError) {
				t.Fail()
			}
		})

	})
}

func BenchmarkParseBase64(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = ParseBase64("7//////////")
	}
}
