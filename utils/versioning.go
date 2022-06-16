package utils

import (
	"strconv"
	"strings"
)

func Version_StrToUint(v string) uint {
	var version [3]string = [3]string{"0", "0", "0"}
	var result uint = 0

	v = strings.Replace(v, "v", "", 1)
	tmpversion := strings.Split(v, ".")

	for i, v_str := range tmpversion {
		version[i] = v_str
	}

	if n, err := strconv.Atoi(version[0]); err != nil {
		return 0
	} else {
		result |= uint(n) << 24
	}

	if n, err := strconv.Atoi(version[1]); err != nil {
		return 0
	} else {
		result |= uint(n) << 16
	}

	if n, err := strconv.Atoi(version[2]); err != nil {
		return 0
	} else {
		result |= uint(n)
	}

	return result
}

func Version_UintToStr(n uint) string {
	var result string
	return result
}
