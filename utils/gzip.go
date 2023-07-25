package utils

import (
	"bytes"
	"compress/gzip"
	"encoding/hex"
)

func GetGzipEncodeStr(str string) (string, error) {
	// 使用gzip对字符串进行压缩
	data, _ := hex.DecodeString(str)
	var compressed bytes.Buffer
	gzipWriter := gzip.NewWriter(&compressed)
	_, err := gzipWriter.Write(data)
	if err != nil {
		return str, err
	}
	gzipWriter.Close()

	compressedStr := compressed.String()
	return compressedStr, nil
}

func GetGzipDecodeStr(str string) (string, error) {
	compressedBytes := []byte(str)
	// 解压缩字符串
	gzipReader, err := gzip.NewReader(bytes.NewReader(compressedBytes))
	if err != nil {
		return str, err
	}
	uncompressedBytes := new(bytes.Buffer)
	_, err = uncompressedBytes.ReadFrom(gzipReader)
	if err != nil {
		return str, err
	}
	return hex.EncodeToString(uncompressedBytes.Bytes()), nil
}
