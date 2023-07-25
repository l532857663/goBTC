package ord

// 铭文唯一ID
func GetInscribeIdStr(txId string) string {
	return txId + "i0"
}

// 铭文输出位置
func GetInscribeOutputStr(txId string) string {
	return txId + ":0"
}
