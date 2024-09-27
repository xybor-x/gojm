package gojm

import "time"

// JobResult is the returned value of Job.
type JobResult struct {
	Err    error
	DoneAt time.Time

	values map[any]any
}

// EmptyResult initializes an empty JobResult.
func EmptyResult() *JobResult {
	return &JobResult{
		Err:    nil,
		values: make(map[any]any),
	}
}

// Result returns a JobResult with a default value.
func Result(value any) *JobResult {
	result := EmptyResult()
	result.Set(nil, value)

	return result
}

// Err return a JobResult with an error.
func Err(err error) *JobResult {
	result := EmptyResult()
	result.Err = err
	return result
}

// Set adds a key-value pair to JobResult.
func (r *JobResult) Set(key, value any) *JobResult {
	r.values[key] = value
	return r
}

// Has returns true if the result contains a key.
func (r JobResult) Has(key any) bool {
	_, ok := r.values[key]
	return ok
}

// Get returns the value with a given key.
func (r JobResult) Get(key any) any {
	return r.values[key]
}

// GetBool returns the boolean value with a given key.
func (r JobResult) GetBool(key any) bool {
	return r.values[key].(bool)
}

// GetInt returns the int value with a given key.
func (r JobResult) GetInt(key any) int {
	return r.values[key].(int)
}

// GetInt32 returns the int32 value with a given key.
func (r JobResult) GetInt32(key any) int32 {
	return r.values[key].(int32)
}

// GetInt64 returns the int64 value with a given key.
func (r JobResult) GetInt64(key any) int64 {
	return r.values[key].(int64)
}

// GetUint returns the uint value with a given key.
func (r JobResult) GetUint(key any) uint {
	return r.values[key].(uint)
}

// GetUint32 returns the uint32 value with a given key.
func (r JobResult) GetUint32(key any) uint32 {
	return r.values[key].(uint32)
}

// GetUint64 returns the uint64 value with a given key.
func (r JobResult) GetUint64(key any) uint64 {
	return r.values[key].(uint64)
}

// GetFloat32 returns the float32 value with a given key.
func (r JobResult) GetFloat32(key any) float32 {
	return r.values[key].(float32)
}

// GetFloat64 returns the float64 value with a given key.
func (r JobResult) GetFloat64(key any) float64 {
	return r.values[key].(float64)
}

// GetString returns the string value with a given key.
func (r JobResult) GetString(key any) string {
	return r.values[key].(string)
}
