package client

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"goBTC/models"
	"strings"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/wire"
)

// 解析原始交易数据
func GetTxWitnessByTxHex(txHex string) string {
	// 将原始十六进制数据解析为交易结构体
	txBytes, err := hex.DecodeString(txHex)
	if err != nil {
		return ""
	}
	var tx wire.MsgTx
	tx.Deserialize(bytes.NewReader(txBytes))
	return GetTxWitness(&tx)
}

// 解析原始交易
func GetTxWitnessByTxInfo(txInfo *btcjson.TxRawResult) string {
	for _, vin := range txInfo.Vin {
		if vin.IsCoinBase() || !vin.HasWitness() {
			continue
		}
		if len(vin.Witness) > 2 {
			return vin.Witness[1]
		}
	}
	return ""
}

func GetTxWitness(tx *wire.MsgTx) string {
	// 遍历交易的输入，查找包含WITNESS_V1_TAPROOT数据的输入
	for _, input := range tx.TxIn {
		// 判断输入是否包含WITNESS_V1_TAPROOT数据
		if len(input.Witness) > 2 {
			// fmt.Printf("wch---- input len(%v): %x\n", len(input.SignatureScript), input.SignatureScript)
			// for _, data := range input.Witness {
			// 	fmt.Printf("wch---- data: %x\n", data)
			// }
			return fmt.Sprintf("%x", input.Witness[1])
		}
	}
	return ""
}

// 解析Input Script inscribe
func GetScriptString(data string) *models.OrdInscribeData {
	char, err := hex.DecodeString(data)
	if err != nil || len(char) == 0 {
		fmt.Println("The string not hex string:", err)
		return nil
	}
	var (
		resObj   = &models.OrdInscribeData{}
		res      = ""
		i        = 0
		headChar = char[0]
		dataType string //铭文类型字符串
	)
	// 初步认定格式
	if headChar == 0x20 {
		// 头部数据公钥
		res, i = getHexData(char, i, "B")
		resObj.PubKey = res
	} else {
		// fmt.Printf("HeadChar: %x\n", headChar)
		return nil
	}
	// 中段数据
	if len(char) < i+3 {
		return nil
	}
	flage1 := hex.EncodeToString(char[i : i+3])
	if flage1 == "ac0063" {
		i += 3 // 跳过ac0063(OP_CHECKSIG OP_FALSE OP_IF)
	} else if char[i] == BTCScriptSignMap["OP_CHECKSIG"] && len(char) > i+3 {
		i += 1 // 跳过ac
		i = skipNothingData(char, i)
	} else {
		fmt.Printf("flage1: [%v] char: [%x]\n", flage1, char)
		return nil
	}
	// Get 铭文 ord [OP_1指示下一次推送包含内容类型，OP_0 指示后续数据推送包含内容本身。]
	res, i = getHexData(char, i, "S")
	if res != "ord" {
		fmt.Printf("%s: %s\n", res, dataType)
		return nil
	}
	// Get 铭文类型
	if len(char) < i+2 {
		return nil
	}
	flage2 := hex.EncodeToString(char[i : i+2])
	if flage2 != "0101" {
		return nil
	}
	i += 2 // 跳过 OP_1
	dataType, i = getHexData(char, i, "S")
	resObj.ContentType = dataType
	tList := strings.Split(dataType, "/")
	inscribeType := tList[0]
	IType, ok := InscribeTypeMap[inscribeType]
	if !ok {
		IType = "B"
	}

	// 后段数据
	if char[i] != 0x00 {
		return nil
	}
	i++ // 跳过OP_0
	lenChar := len(char[i:])
	res = ""
	start := 0
	// fmt.Printf("filer lenChar: %+v, char[i]: %x, i: %+v\n", lenChar, char[i], i)
	for j := 0; j <= lenChar; j += len(char[start:i]) {
		start = i
		tmp := ""
		if char[i] >= 0x01 && char[i] <= 0x4b {
			tmp, i = getHexData(char, i, IType)
		} else if char[i] == 0x4c {
			i++ // OP_PUSHDATA1
			tmp, i = getHexData(char, i, IType)
		} else if char[i] == 0x4d {
			hexLen := getHexLen(char[i+1 : i+3])
			i += 3 // OP_PUSHDATA2
			tmp, i = getHexPushData(char, hexLen, i, IType)
		} else if char[i] == 0x4e {
			hexLen := getHexLen(char[i+1 : i+5])
			i += 5 // OP_PUSHDATA4
			tmp, i = getHexPushData(char, hexLen, i, IType)
		} else {
			opCode := char[i]
			if opCode == 0x68 {
				fmt.Println("End!")
				break
			}
			fmt.Printf("End?: %x, %s, [%s]\n", opCode, string(opCode), BTCScriptValueMap[opCode])
			break
		}
		res += tmp
	}
	// 获取铭文数据
	resObj.Body = res
	resLength := int64(len(res))
	resObj.ContentSize = resLength / 2
	if inscribeType == InscribeTypeText || inscribeType == InscribeTypeApplication {
		json.Unmarshal([]byte(res), &resObj.Brc20)
		if resObj.Brc20 == nil {
			resB, err := hex.DecodeString(res)
			if err == nil {
				json.Unmarshal(resB, &resObj.Brc20)
				if resObj.Brc20 != nil {
					resObj.Body = string(resB)
				}
			}
		}
		resObj.ContentSize = resLength
	}
	return resObj
}

// 获取数据内容
func getHexData(char []byte, i int, getType string) (string, int) {
	hexLen := int(char[i])
	i++
	return getHexPushData(char, hexLen, i, getType)
}

// 根据长度获取数据内容
func getHexPushData(char []byte, hexLen, i int, getType string) (string, int) {
	data := char[i : i+hexLen]
	i += hexLen
	res := ""
	switch getType {
	case "B":
		res = hex.EncodeToString(data)
	case "S":
		res = string(data)
	}
	return res, i
}

// 跳过无用数据
func skipNothingData(char []byte, i int) int {
	_, i = getHexData(char, i, "S")
	if char[i] >= 0x01 && char[i] <= 0x4b {
		i = skipNothingData(char, i)
	} else {
		i++
		flage := hex.EncodeToString(char[i : i+2])
		if flage == "0063" {
			i += 2 // 跳过0063(OP_FALSE OP_IF)
		}
	}
	return i
}

func getHexLen(char []byte) int {
	var num uint32
	for i := len(char) - 1; i >= 0; i-- {
		num = num<<8 + uint32(char[i]) // 将字节数组转换为数字
	}
	return int(num)
}

func GetWitnessScript(data string) []string {
	char, err := hex.DecodeString(data)
	if err != nil {
		fmt.Println("The string not hex string")
		return nil
	}
	for _, c := range char {
		fmt.Printf("char: %x   %s   %s\n", c, string(c), BTCScriptValueMap[c])
	}
	return nil
}
