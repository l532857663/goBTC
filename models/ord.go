package models

type OrdInscribeData struct {
	PubKey      string
	ContentType string
	ContentSize int64
	Body        string
	Destination string
	Brc20       *OrdBRC20
	TxFee       int64
	// 该交易是否是铭文UTXO转账
	TxHaveInscribe bool
}

type OrdBRC20 struct {
	P    string `json:"p"`
	OP   string `json:"op"`
	Tick string `json:"tick"`
	Amt  string `json:"amt"`
	Lim  string `json:"lim"`
	Max  string `json:"max"`
}
