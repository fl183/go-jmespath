package jputil

import (
	"errors"
	"reflect"
)

// IsFalse determines if an object is false based on the JMESPath spec.
// JMESPath defines false values to be any of:
// - An empty string array, or hash.
// - The boolean value false.
// - nil
func IsFalse(value interface{}) bool {
	if value == nil {
		return true
	} else if value == false {
		return true
	} else if aSlice, ok := value.([]interface{}); ok && len(aSlice) == 0 {
		return true
	} else if aMap, ok := value.(map[string]interface{}); ok && len(aMap) == 0 {
		return true
	} else if aStr, ok := value.(string); ok && len(aStr) == 0 {
		return true
	}
	return false
}

// ObjsEqual is a generic object equality check.
// It will take two arbitrary objects and recursively determine
// if they are equal.
func ObjsEqual(left interface{}, right interface{}) bool {
	if (left == nil) || (right == nil) {
		return left == right
	}
	if reflect.DeepEqual(left, right) {
		return true
	}
	return false
}

// SliceParam refers to a single part of a slice.
// A slice consists of a start, a stop, and a step, similar to
// python slices.
type SliceParam struct {
	N         int
	Specified bool
}

// Slice supports [start:stop:step] style slicing that's supported in JMESPath.
func Slice(slice []interface{}, parts []SliceParam) ([]interface{}, error) {
	computed, err := computeSliceParams(len(slice), parts)
	if err != nil {
		return nil, err
	}
	start, stop, step := computed[0], computed[1], computed[2]
	result := make([]interface{}, 0, 0)
	if step > 0 {
		for i := start; i < stop; i += step {
			result = append(result, slice[i])
		}
	} else {
		for i := start; i > stop; i += step {
			result = append(result, slice[i])
		}
	}
	return result, nil
}

func computeSliceParams(length int, parts []SliceParam) ([]int, error) {
	var start, stop, step int
	if !parts[2].Specified {
		step = 1
	} else if parts[2].N == 0 {
		return nil, errors.New("Invalid slice, step cannot be 0")
	} else {
		step = parts[2].N
	}
	var stepValueNegative bool
	if step < 0 {
		stepValueNegative = true
	} else {
		stepValueNegative = false
	}

	if !parts[0].Specified {
		if stepValueNegative {
			start = length - 1
		} else {
			start = 0
		}
	} else {
		start = capSlice(length, parts[0].N, step)
	}

	if !parts[1].Specified {
		if stepValueNegative {
			stop = -1
		} else {
			stop = length
		}
	} else {
		stop = capSlice(length, parts[1].N, step)
	}
	return []int{start, stop, step}, nil
}

func capSlice(length int, actual int, step int) int {
	if actual < 0 {
		actual += length
		if actual < 0 {
			if step < 0 {
				actual = -1
			} else {
				actual = 0
			}
		}
	} else if actual >= length {
		if step < 0 {
			actual = length - 1
		} else {
			actual = length
		}
	}
	return actual
}

// ToArrayNum converts an empty interface type to a slice of float64.
// If any element in the array cannot be converted, then nil is returned
// along with a second value of false.
func ToArrayNum(data interface{}) ([]float64, bool) {
	// Is there a better way to do this with reflect?
	if d, ok := data.([]interface{}); ok {
		result := make([]float64, len(d))
		for i, el := range d {
			item, ok := el.(float64)
			if !ok {
				return nil, false
			}
			result[i] = item
		}
		return result, true
	}
	return nil, false
}

// ToArrayStr converts an empty interface type to a slice of strings.
// If any element in the array cannot be converted, then nil is returned
// along with a second value of false.  If the input data could be entirely
// converted, then the converted data, along with a second value of true,
// will be returned.
func ToArrayStr(data interface{}) ([]string, bool) {
	// Is there a better way to do this with reflect?
	if d, ok := data.([]interface{}); ok {
		result := make([]string, len(d))
		for i, el := range d {
			item, ok := el.(string)
			if !ok {
				return nil, false
			}
			result[i] = item
		}
		return result, true
	}
	return nil, false
}