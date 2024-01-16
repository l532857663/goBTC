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
	ContentType       string
	Body              string
	Destination       string
	TxFee             int64
	ChangeAddress     string
	ServiceFeeAddress string
	ServiceFee        int64
}

type OrdiInfo struct {
	Tick        string
	Amount      int64
	Body        string
	To          string
	ServiceFee  int64
	GasFee      int64
	GasFeeTotal int64
	PSBTData    string
	ContentType string
}
