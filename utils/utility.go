package utils

import "strconv"

func IntToHex(n int64) []byte { // 10진수 정수를 16진법으로 변경
	return []byte(strconv.FormatInt(n, 16))
}
