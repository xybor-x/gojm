package gojm_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xybor-x/gojm"
)

func Test_JobResult_Has(t *testing.T) {
	result := gojm.EmptyResult()

	assert.False(t, result.Has("var1"))

	result.Set("var1", 3)
	assert.True(t, result.Has("var1"))
}

func Test_JobResult_Get(t *testing.T) {
	result := gojm.EmptyResult()

	result.Set("var1", 3)
	assert.Equal(t, 3, result.Get("var1"))
}

func Test_JobResult_Error(t *testing.T) {
	expectedErr := errors.New("something wrong")

	result := gojm.Err(expectedErr)

	assert.ErrorIs(t, result.Err, expectedErr)
}

func Test_JobResult_GetType(t *testing.T) {
	result := gojm.EmptyResult()

	result.Set("bool", true)
	assert.Equal(t, true, result.GetBool("bool"))

	result.Set("int", 1)
	assert.Equal(t, 1, result.GetInt("int"))

	result.Set("int32", int32(1))
	assert.Equal(t, int32(1), result.GetInt32("int32"))

	result.Set("int64", int64(1))
	assert.Equal(t, int64(1), result.GetInt64("int64"))

	result.Set("uint", uint(1))
	assert.Equal(t, uint(1), result.GetUint("uint"))

	result.Set("uint32", uint32(1))
	assert.Equal(t, uint32(1), result.GetUint32("uint32"))

	result.Set("uint64", uint64(1))
	assert.Equal(t, uint64(1), result.GetUint64("uint64"))

	result.Set("float32", float32(1))
	assert.Equal(t, float32(1), result.GetFloat32("float32"))

	result.Set("float64", float64(1))
	assert.Equal(t, float64(1), result.GetFloat64("float64"))

	result.Set("string", "something")
	assert.Equal(t, "something", result.GetString("string"))
}
