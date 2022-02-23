package database_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"lexibot/internal/database"
	"testing"
)

func TestStringArray_Scan(t *testing.T) {
	t.Run("scan string", func(t *testing.T) {
		src := database.StringArray{}
		err := src.Scan("first,second")

		require.NoError(t, err)
		assert.Equal(t, database.StringArray{"first", "second"}, src)
	})

	t.Run("scan number", func(t *testing.T) {
		src := database.StringArray{}
		err := src.Scan(10)

		require.Error(t, err)
	})
}

func TestStringArray_Value(t *testing.T) {
	for _, src := range []database.StringArray{{}, nil} {
		t.Run("empty value", func(t *testing.T) {
			value, err := src.Value()

			require.NoError(t, err)
			assert.Nil(t, value)
		})
	}

	t.Run("string value", func(t *testing.T) {
		src := database.StringArray{"first", "second"}
		value, err := src.Value()

		require.NoError(t, err)
		assert.Equal(t, "first,second", value)
	})
}

func TestStringArray_GormDataType(t *testing.T) {
	assert.Equal(t, "text", database.StringArray{}.GormDataType())
}
