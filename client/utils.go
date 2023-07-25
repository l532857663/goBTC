package client

import (
	"fmt"
)

func StuffMapBySlice(signMap map[string]byte, valueMap map[byte]string, signSlice []string, valSlice []byte) {
	if len(signSlice) != len(valSlice) {
		fmt.Printf("StuffMapBySlice lenSign: %+v, lenValue: %+v\n", len(signSlice), len(valSlice))
		return
	}
	for i, sign := range signSlice {
		val := valSlice[i]
		signMap[sign] = val
		valueMap[val] = sign
	}
	return
}
