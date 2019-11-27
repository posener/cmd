package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestArgs(t *testing.T) {
	t.Parallel()

	t.Run("no cap", func(t *testing.T) {
		var args ArgsStr
		err := args.Set([]string{"a", "b"})
		require.NoError(t, err)
		assert.Equal(t, []string{"a", "b"}, []string(args))
	})

	t.Run("with cap", func(t *testing.T) {
		var args = make(ArgsStr, 2)
		err := args.Set([]string{"a", "b"})
		require.NoError(t, err)
		assert.Equal(t, []string{"a", "b"}, []string(args))
	})

	t.Run("too many args", func(t *testing.T) {
		var args = make(ArgsStr, 2)
		err := args.Set([]string{"a", "b", "c"})
		require.Error(t, err)
	})

	t.Run("not enough args", func(t *testing.T) {
		var args = make(ArgsStr, 2)
		err := args.Set([]string{"a"})
		require.Error(t, err)
	})
}

func TestArgsInt(t *testing.T) {
	t.Parallel()

	t.Run("no cap", func(t *testing.T) {
		var args ArgsInt
		err := args.Set([]string{"1", "2"})
		require.NoError(t, err)
		assert.Equal(t, []int{1, 2}, []int(args))
	})

	t.Run("with cap", func(t *testing.T) {
		var args = make(ArgsInt, 2)
		err := args.Set([]string{"1", "2"})
		require.NoError(t, err)
		assert.Equal(t, []int{1, 2}, []int(args))
	})

	t.Run("bad value", func(t *testing.T) {
		var args ArgsInt
		err := args.Set([]string{"a"})
		require.Error(t, err)
	})

	t.Run("too many args", func(t *testing.T) {
		var args = make(ArgsInt, 2)
		err := args.Set([]string{"1", "2", "3"})
		require.Error(t, err)
	})

	t.Run("not enough args", func(t *testing.T) {
		var args = make(ArgsInt, 2)
		err := args.Set([]string{"1"})
		require.Error(t, err)
	})
}
