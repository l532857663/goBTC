package client

// 常用节点数据 NodeInfo ----------------------------------------------------------
const (
	// 节点类型
	MainNet = "MAINNET"
	TestNet = "TESTNET3"
	RegNet  = "REGNET"
)

// NodeInfo ----------------------------------------------------------

// BTC Script Info ----------------------------------------------------------
var (
	// 全局变量
	BTCScriptSignMap  = make(map[string]byte)
	BTCScriptValueMap = make(map[byte]string)

	// 脚本压入堆栈
	Sign1  = []string{"OP_0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16", "OP_PUSHDATA1", "OP_PUSHDATA2", "OP_PUSHDATA4", "OP_1NEGATE", "OP_RESERVED", "OP_1", "OP_2", "OP_3", "OP_4", "OP_5", "OP_6", "OP_7", "OP_8", "OP_9", "OP_10", "OP_11", "OP_12", "OP_13", "OP_14", "OP_15", "OP_16"}
	Value1 = []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10, 0x4c, 0x4d, 0x4e, 0x4f, 0x50, 0x51, 0x52, 0x53, 0x54, 0x55, 0x56, 0x57, 0x58, 0x59, 0x5a, 0x5b, 0x5c, 0x5d, 0x5e, 0x5f, 0x60}
	// 有条件的流控制的操作符
	Sign2  = []string{"OP_NOP", "OP_VER", "OP_IF", "OP_NOTIF", "OP_VERIF", "OP_VERNOTIF", "OP_ELSE", "OP_ENDIF", "OP_VERIFY", "OP_RETURN"}
	Value2 = []byte{0x61, 0x62, 0x63, 0x64, 0x65, 0x66, 0x67, 0x68, 0x69, 0x6a}
	// 时间锁操作符
	Sign3  = []string{"OP_CHECKLOCKTIMEVERIFY", "OP_CHECKSEQUENCEVERIFY"}
	Value3 = []byte{0xb1, 0xb2}
	// 堆栈操作符
	Sign4  = []string{"OP_TOALTSTACK", "OP_FROMALTSTACK", "OP_2DROP", "OP_2DUP", "OP_3DUP", "OP_2OVER", "OP_2ROT", "OP_2SWAP", "OP_IFDUP", "OP_DEPTH", "OP_DROP", "OP_DUP", "OP_NIP", "OP_OVER", "OP_PICK", "OP_ROLL", "OP_ROT", "OP_SWAP", "OP_TUCK"}
	Value4 = []byte{0x6b, 0x6c, 0x6d, 0x6e, 0x6f, 0x70, 0x71, 0x72, 0x73, 0x74, 0x75, 0x76, 0x77, 0x78, 0x79, 0x7a, 0x7b, 0x7c, 0x7d}
	// 字符串操作符
	Sign5  = []string{"OP_CAT", "OP_SUBSTR", "OP_LEFT", "OP_RIGHT", "OP_SIZE"}
	Value5 = []byte{0x7e, 0x7f, 0x80, 0x81, 0x82}
	// 二进制算术和条件
	Sign6  = []string{"OP_INVERT", "OP_AND", "OP_OR", "OP_XOR", "OP_EQUAL", "OP_EQUALVERIFY", "OP_RESERVED1", "OP_RESERVED2"}
	Value6 = []byte{0x83, 0x84, 0x85, 0x86, 0x87, 0x88, 0x89, 0x8a}
	// 数值操作
	Sign7 = []string{"OP_1ADD", "OP_1SUB", "OP_2MUL", "OP_2DIV", "OP_NEGATE", "OP_ABS", "OP_NOT", "OP_0NOTEQUAL", "OP_ADD", "OP_SUB", "OP_MUL", "OP_DIV", "OP_MOD", "OP_LSHIFT", "OP_RSHIFT", "OP_BOOLAND", "OP_BOOLOR", "OP_NUMEQUAL", "OP_NUMEQUALVERIFY", "OP_NUMNOTEQUAL", "OP_LESSTHAN", "OP_GREATERTHAN", "OP_LESSTHANOREQUAL", "OP_GREATERTHANOREQUAL", "OP_MIN", "OP_MAX", "OP_WITHIN"}

	Value7 = []byte{0x8b, 0x8c, 0x8d, 0x8e, 0x8f, 0x90, 0x91, 0x92, 0x93, 0x94, 0x95, 0x96, 0x97, 0x98, 0x99, 0x9a, 0x9b, 0x9c, 0x9d, 0x9e, 0x9f, 0xa0, 0xa1, 0xa2, 0xa3, 0xa4, 0xa5}
	// 加密和散列操作
	Sign8  = []string{"OP_RIPEMD160", "OP_SHA1", "OP_SHA256", "OP_HASH160", "OP_HASH256", "OP_CODESEPARATOR", "OP_CHECKSIG", "OP_CHECKSIGVERIFY", "OP_CHECKMULTISIG", "OP_CHECKMULTISIGVERIFY"}
	Value8 = []byte{0xa6, 0xa7, 0xa8, 0xa9, 0xaa, 0xab, 0xac, 0xad, 0xae, 0xaf}
	// 仅供内部使用的保留关键字
	Sign9  = []string{"OP_SMALLDATA", "OP_SMALLINTEGER", "OP_PUBKEYS", "OP_PUBKEYHASH", "OP_PUBKEY", "OP_INVALIDOPCODE"}
	Value9 = []byte{0xf9, 0xfa, 0xfb, 0xfd, 0xfe, 0xff}
	// 字符串接操作都已禁用(除0x82)
	Sign10  = []string{"OP_CAT", "OP_SUBSTR", "OP_LEFT", "OP_RIGHT", "OP_SIZE"}
	Value10 = []byte{0x7e, 0x7f, 0x80, 0x81, 0x82}
)

func InitBtcScriptMap() {
	StuffMapBySlice(BTCScriptSignMap, BTCScriptValueMap, Sign1, Value1)
	StuffMapBySlice(BTCScriptSignMap, BTCScriptValueMap, Sign2, Value2)
	StuffMapBySlice(BTCScriptSignMap, BTCScriptValueMap, Sign3, Value3)
	StuffMapBySlice(BTCScriptSignMap, BTCScriptValueMap, Sign4, Value4)
	StuffMapBySlice(BTCScriptSignMap, BTCScriptValueMap, Sign5, Value5)
	StuffMapBySlice(BTCScriptSignMap, BTCScriptValueMap, Sign6, Value6)
	StuffMapBySlice(BTCScriptSignMap, BTCScriptValueMap, Sign7, Value7)
	StuffMapBySlice(BTCScriptSignMap, BTCScriptValueMap, Sign8, Value8)
	StuffMapBySlice(BTCScriptSignMap, BTCScriptValueMap, Sign9, Value9)
	StuffMapBySlice(BTCScriptSignMap, BTCScriptValueMap, Sign10, Value10)
}

// BTC Script Info ----------------------------------------------------------

// BTC Script InscribeType --------------------------------------------------
const (
	InscribeTypeText        = "text"
	InscribeTypeImage       = "image"
	InscribeTypeVideo       = "video"
	InscribeTypeAudio       = "audio"
	InscribeTypeApplication = "application"
)

var (
	InscribeTypeMap = map[string]string{
		InscribeTypeText:        "S",
		InscribeTypeImage:       "B",
		InscribeTypeVideo:       "B",
		InscribeTypeAudio:       "B",
		InscribeTypeApplication: "B",
	}
)

// BTC Script InscribeType --------------------------------------------------
