package utils

import "strconv"

func StringToUint16(textValue string) (uint16, error) {
	uint64Value, err := strconv.ParseUint(textValue, 10, 64)
	if err != nil {
		return 0, err
	}

	return uint16(uint64Value), nil
}
