package toolbox_test

import (
	"testing"
	"github.com/viant/toolbox"
	"github.com/stretchr/testify/assert"
)

func TestSliceIterator(t *testing.T) {

	{
		slice := []string{"a", "r", "c"}
		iterator := toolbox.NewSliceIterator(slice)
		var values = make([]interface{}, 1)
		value := values[0]

		assert.True(t, iterator.HasNext())
		iterator.Next(&value)
		assert.Equal(t, "a", value)

		assert.True(t, iterator.HasNext())
		iterator.Next(&value)
		assert.Equal(t, "r", value)



		assert.True(t, iterator.HasNext())
		iterator.Next(&value)

		assert.Equal(t, "c", value)

	}
	{
		slice := []string{"a", "r", "c"}
		iterator := toolbox.NewSliceIterator(slice)
		value := ""

		assert.True(t, iterator.HasNext())
		iterator.Next(&value)
		assert.Equal(t, "a", value)

		assert.True(t, iterator.HasNext())
		iterator.Next(&value)
		assert.Equal(t, "r", value)



		assert.True(t, iterator.HasNext())
		iterator.Next(&value)

		assert.Equal(t, "c", value)

	}
	{
		slice := []interface{}{"a", "z", "c"}
		iterator := toolbox.NewSliceIterator(slice)
		value := ""

		assert.True(t, iterator.HasNext())
		iterator.Next(&value)
		assert.Equal(t, "a", value)

		assert.True(t, iterator.HasNext())
		iterator.Next(&value)
		assert.Equal(t, "z", value)


		var values = make([]interface{}, 1)
		assert.True(t, iterator.HasNext())
		iterator.Next(&values[0])
		assert.Equal(t, "c", values[0])

	}

	{
		slice := []int{3, 2, 1}
		iterator := toolbox.NewSliceIterator(slice)
		value := 0
		assert.True(t, iterator.HasNext())
		iterator.Next(&value)
		assert.Equal(t, 3, value)

		assert.True(t, iterator.HasNext())
		iterator.Next(&value)
		assert.Equal(t, 2, value)

		assert.True(t, iterator.HasNext())
		iterator.Next(&value)
		assert.Equal(t, 1, value)
	}

}

