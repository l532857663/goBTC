package models

type PageInfo struct {
	Page  string `json:"page"`
	Limit string `json:"limit"`
}

type Filter struct {
	PageInfo
	ContractAddress string
	TransferType    string
}

type InscribeFilter struct {
	PageInfo
	Token             string
	InscriptionId     string
	InscriptionNumber string
	State             string
}

type GetBalanceParam struct {
	PageInfo
	InscriptionId string `json:"inscription_id"`
	Address       string `json:"address"`
}

type CreateTransferParam struct {
	Address string `json:"address"`
	Tick    string `json:"tick"`
	Amount  string `json:"amount"`
	FeeRate int64  `json:"feeRate"`
}

type SendTransferParam struct {
	PSBTData string `json:"psbtData"`
	Key      string `json:"key"`
}
