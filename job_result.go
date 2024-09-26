package gojm

import "time"

type JobResult struct {
	Err    error
	DoneAt time.Time

	values map[any]any
}

func EmptyResult() *JobResult {
	return &JobResult{
		Err:    nil,
		values: make(map[any]any),
	}
}

func Result(value any) *JobResult {
	result := EmptyResult()
	result.Set(nil, value)

	return result
}

func Err(err error) *JobResult {
	result := EmptyResult()
	result.Err = err
	return result
}

func (r *JobResult) Set(key, value any) *JobResult {
	r.values[key] = value
	return r
}

func (r JobResult) Has(key any) bool {
	_, ok := r.values[key]
	return ok
}

func (r JobResult) Get(key any) any {
	return r.values[key]
}

func (r JobResult) GetBool(key any) bool {
	return r.values[key].(bool)
}

func (r JobResult) GetInt(key any) int {
	return r.values[key].(int)
}

func (r JobResult) GetInt32(key any) int32 {
	return r.values[key].(int32)
}

func (r JobResult) GetInt64(key any) int64 {
	return r.values[key].(int64)
}

func (r JobResult) GetUint(key any) uint {
	return r.values[key].(uint)
}

func (r JobResult) GetUint32(key any) uint32 {
	return r.values[key].(uint32)
}

func (r JobResult) GetUint64(key any) uint64 {
	return r.values[key].(uint64)
}

func (r JobResult) GetFloat32(key any) float32 {
	return r.values[key].(float32)
}

func (r JobResult) GetFloat64(key any) float64 {
	return r.values[key].(float64)
}

func (r JobResult) GetString(key any) string {
	return r.values[key].(string)
}
