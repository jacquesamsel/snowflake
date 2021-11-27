package snowflake

import (
	"errors"
	"github.com/bmizerany/assert"
	"testing"
	"time"
)

func TestNewNode(t *testing.T) {
	_, err := NewNode(0, time.Now(), 42, 5, 16)
	assert.Equal(t, nil, err)
	_, err = NewNode(0, time.Now(), 43, 5, 16)
	assert.T(t, errors.Is(err, ErrorSnowflakeOverflow))
	_, err = NewNode(33, time.Now(), 42, 5, 16)
	assert.T(t, errors.Is(err, ErrorNodeOverflow))
}

func TestNode_Generate(t *testing.T) {
	node, err := NewNode(0, time.Now(), 42, 2, 19)
	assert.Equal(t, nil, err)
	t.Run("single", func(t *testing.T) {
		id := node.Generate()
		assert.NotEqual(t, id, 0)
	})
	t.Run("uniqueness", func(t *testing.T) {
		var x, y Snowflake
		for i := 0; i < 1e6; i++ {
			x, y = y, node.Generate()
			assert.NotEqual(t, x, y)
		}
	})
}

func BenchmarkNode_Generate(b *testing.B) {
	b.ReportAllocs()
	node, err := NewNode(0, time.Now().Add(-10*time.Hour), 42, 1, 20)
	if err != nil {
		panic(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		node.Generate()
	}
}
