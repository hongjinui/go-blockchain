package utils

import "strconv"

func IntToHex(n int64) []byte { // 10진수 정수를 16진법으로 변경
	return []byte(strconv.FormatInt(n, 16))

}

// ReverseBtyes reverse a byte array
func ReverseBytes(data []byte) {
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
}
