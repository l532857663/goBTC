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
	TxHaveInscribe string
}

type OrdBRC20 struct {
	P    string `json:"p"`
	OP   string `json:"op,omitempty"`
	Tick string `json:"tick,omitempty"`
	Amt  string `json:"amt,omitempty"`
	Lim  string `json:"lim,omitempty"`
	Max  string `json:"max,omitempty"`
	Name string `json:"name,omitempty"`
}

type CreateOrdFilter struct {
	ContentType   string
	Body          string
	Destination   string
	TxFee         int64
	ChangeAddress string
}
