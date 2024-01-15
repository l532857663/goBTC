package models

type GetBalanceResp struct {
	Balance        string      `json:"balance"`
	BalanceSymbol  string      `json:"balanceSymbol"`
	ChainFullName  string      `json:"chainFullName"`
	ChainShortName string      `json:"chainShortName"`
	TotalPage      string      `json:"totalPage"`
	TokenList      []TokenInfo `json:"tokenList"`
}

type TokenInfo struct {
	Token                string `json:"token"`
	TokenContractAddress string `json:"tokenContractAddress"`
	Amount               string `json:"amount"`
}

type GetTransferResp struct {
	ChainFullName  string         `json:"chainFullName"`
	ChainShortName string         `json:"chainShortName"`
	TotalPage      string         `json:"totalPage"`
	TransferList   []TransferInfo `json:"transferList"`
}

type TransferInfo struct {
	Txid                 string `json:"txId"`                 // 交易哈希
	BlockHash            string `json:"blockHash"`            // 区块哈希
	Height               string `json:"height"`               // 交易发生的区块
	TransactionTime      string `json:"transactionTime"`      // 交易时间；Unix时间戳的毫秒数格式，如 1597026383085
	From                 string `json:"from"`                 // 发送方地址
	To                   string `json:"to"`                   // 接收方地址
	IsFromContract       bool   `json:"isFromContract"`       // From地址是否是合约地址
	IsToContract         bool   `json:"isToContract"`         // To地址是否是合约地址
	FromTag              string `json:"fromTag"`              // 发送方地址标签
	ToTag                string `json:"toTag"`                // 接收方地址标签
	Amount               string `json:"amount"`               // 交易数量
	TransactionSymbol    string `json:"transactionSymbol"`    // 交易数量对应的币种
	TokenContractAddress string `json:"tokenContractAddress"` // 交易数量对应的币种的合约地址
	TxFee                string `json:"txfee"`                // 手续费
	State                string `json:"state"`                // 交易状态 success 成功 fail 失败 pending 等待确认
	// ByHeight
	TokenID  string `json:"tokenId"`  // NFT的ID，如果该交易是NFT交易
	MethodId string `json:"methodId"` // 方法，如果该交易是合约调用交易
}

type GetTransferUTXOResp struct {
	ChainFullName    string `json:"chainFullName"`    // 公链全称，例如：Bitcoin
	ChainShortName   string `json:"chainShortName"`   // 公链缩写符号，例如：BTC
	Txid             string `json:"txId"`             // 交易哈希
	Height           string `json:"height"`           // 交易发生的区块
	Amount           string `json:"amount"`           // UTXO里面的交易金额
	Address          string `json:"address"`          // 地址
	Unspent          string `json:"unspent"`          // 该笔交易未花费的交易输出（找零）
	Confirm          string `json:"confirm"`          // 确认数
	Index            string `json:"index"`            // 该笔交易交易在区块里的位置索引
	TransactionIndex string `json:"transactionIndex"` // 该笔UTXO在交易里的位置索引
	Balance          string `json:"balance"`          // 该地址余额
	Symbol           string `json:"symbol"`           // 币种
}

type GetInscribeResp struct {
	Page             string             `json:"page"`
	Limit            string             `json:"limit"`
	TotalPage        string             `json:"totalPage"`
	TotalInscription string             `json:"totalInscription"`
	InscriptionsList []InscriptionsList `json:"inscriptionsList"`
}

type InscriptionsList struct {
	InscriptionId     string `json:"inscriptionId"`
	InscriptionNumber string `json:"inscriptionNumber"`
	Location          string `json:"location"`
	Token             string `json:"token"`
	State             string `json:"state"`
	Msg               string `json:"msg"`
	TokenType         string `json:"tokenType"`
	ActionType        string `json:"actionType"`
	LogoURL           string `json:"logoUrl"`
	OwnerAddress      string `json:"ownerAddress"`
	TxID              string `json:"txId"`
	BlockHeight       string `json:"blockHeight"`
	ContentSize       string `json:"contentSize"`
	Time              string `json:"time"`
}

type CreateTransferResp struct {
	PSBTData    string `json:"psbtData"`
	Key         string `json:"key"`
	CommitFee   int64  `json:"commitFee"`
	RevealFee   int64  `json:"revealFee"`
	RevealValue int64  `json:"revealValue"`
	ServiceFee  string `json:"serviceFee"`
}

type SendTransferResp struct {
	CommitTxHash string   `json:"commitTxHash"`
	RevealTxHash []string `json:"revealTxHash"`
}
