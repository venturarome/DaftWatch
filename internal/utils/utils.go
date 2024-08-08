package utils

import "strconv"

func StringToInt(textValue string) int {
	int64Value, err := strconv.ParseInt(textValue, 10, 64)
	if err != nil {
		return 0
	}

	return int(int64Value)
}

func BoolPtr(b bool) *bool {
	return &b
}

// Use of generics: https://go.dev/tour/generics/1
func SliceToInterfaceSlice[T any](tSlice []T) []interface{} {
	iSlice := make([]interface{}, len(tSlice))

	for i := range tSlice {
		iSlice[i] = tSlice[i]
	}

	return iSlice
}

func MapSlice[T any, U any](slice []T, fn func(T) U) []U {
	result := make([]U, len(slice))
	for i, v := range slice {
		result[i] = fn(v)
	}
	return result
}
