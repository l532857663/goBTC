package ord

import "strings"

// 铭文唯一ID
func GetInscribeIdStr(txId string) string {
	return txId + "i0"
}

// 铭文输出位置
func GetInscribeOutputStr(txId string) string {
	return txId + ":0"
}

// BRC20的ID
func GetDeployIdStr(tick string) string {
	return strings.ToLower(tick)
}
