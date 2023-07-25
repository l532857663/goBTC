package oklink

type BaseFilter struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
}

type BaseData interface {
	*AddressInfo | *AddressTokenInfo | *TransactionInfo | *TransactionUTXO | *InscriptionsInfo
}

type BaseResp[T BaseData] struct {
	BaseFilter
	Data []T
}

// 地址基础信息
type AddressInfo struct {
	ChainFullName                 string `json:"chainFullName"`                 // 公链全称，例如：Bitcoin
	ChainShortName                string `json:"chainShortName"`                // 公链缩写符号，例如：BTC
	Address                       string `json:"address"`                       // 普通地址
	ContractAddress               string `json:"contractAddress"`               // 智能合约地址
	IsProducerAddress             bool   `json:"isProducerAddress"`             // 是否为验证人地址，true：是验证人/节点地址，false：不是验证人/节点地址
	Tag                           string `json:"tag"`                           // 地址标签
	Balance                       string `json:"balance"`                       // 余额
	BalanceSymbol                 string `json:"balanceSymbol"`                 // 余额币种
	TransactionCount              string `json:"transactionCount"`              // 该地址交易次数
	Verifying                     string `json:"verifying"`                     // 确认中金额
	SendAmount                    string `json:"sendAmount"`                    // 发送金额
	ReceiveAmount                 string `json:"receiveAmount"`                 // 接收金额
	TokenAmount                   string `json:"tokenAmount"`                   // 代币种类数量
	TotalTokenValue               string `json:"totalTokenValue"`               // 代币总价值折算成公链原生币的数量
	CreateContractAddress         string `json:"createContractAddress"`         // 创建该智能合约的地址
	CreateContractTransactionHash string `json:"createContractTransactionHash"` // 创建该智能合约的交易hash
	FirstTransactionTime          string `json:"firstTransactionTime"`          // 该地址发生第一笔交易时间
	LastTransactionTime           string `json:"lastTransactionTime"`           // 该地址最近一次发生交易时间
	Token                         string `json:"token"`                         // 该地址对应的代币
	Bandwidth                     string `json:"bandwidth"`                     // 带宽和消耗的带宽（仅适用于TRON）
	Energy                        string `json:"energy"`                        // 能量和消耗的能量（仅适用于TRON），其他链返回“”
	VotingRights                  string `json:"votingRights"`                  // 投票全和已用投票权（仅适用于TRON）
	UnclaimedVotingRewards        string `json:"unclaimedVotingRewards"`        // 待领取投票收益（仅适用于TRON）
}

// 代币信息
type AddressTokenInfo struct {
	Page           string      `json:"page"`           // 当前页码
	Limit          string      `json:"limit"`          // 当前页共多少条数据
	TotalPage      string      `json:"totalPage"`      // 总共多少页
	ChainFullName  string      `json:"chainFullName"`  // 公链全称，例如：Bitcoin
	ChainShortName string      `json:"chainShortName"` // 公链缩写符号，例如：BTC
	TokenList      []TokenInfo `json:"tokenList"`      // 代币列表
}

type TokenInfo struct {
	Token                string `json:"token"`                // 该地址对应的代币
	TokenContractAddress string `json:"tokenContractAddress"` // 该地址对应的代币合约地址
	HoldingAmount        string `json:"holdingAmount"`        // 代币持仓数量
	TotalTokenValue      string `json:"totalTokenValue"`      // 代币总价值折算成公链原生币的数量
	Change24h            string `json:"change24h"`            // 代币价格24小时涨跌幅
	PriceUsd             string `json:"priceUsd"`             // 代币美元价格
	ValueUsd             string `json:"valueUsd"`             // 代币总的美元价值
}

// 交易信息
type TransactionInfo struct {
	Page            string        `json:"page"`             // 当前页码
	Limit           string        `json:"limit"`            // 当前页共多少条数据
	TotalPage       string        `json:"totalPage"`        // 总共多少页
	ChainFullName   string        `json:"chainFullName"`    // 公链全称，例如：Bitcoin
	ChainShortName  string        `json:"chainShortName"`   // 公链缩写符号，例如：BTC
	TransactionList []Transaction `json:"transactionLists"` // 交易列表
	BlockList       []Transaction `json:"blockList"`        // 交易列表
}

type Transaction struct {
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

// 交易信息UTXO
type TransactionUTXOInfo struct {
	BaseFilter
	Data []TransactionUTXO
}
type TransactionUTXO struct {
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

// BTC链的 inscriptions 列表
type InscriptionsInfo struct {
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
