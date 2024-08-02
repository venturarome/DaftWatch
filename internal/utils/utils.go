package utils

import "strconv"

func StringToInt(textValue string) (int, error) {
	int64Value, err := strconv.ParseInt(textValue, 10, 64)
	if err != nil {
		return 0, err
	}

	return int(int64Value), nil
}

func BoolPtr(b bool) *bool {
	return &b
}
