package models

type Filter struct {
	Page            string
	Limit           string
	ContractAddress string
	TransferType    string
}

type InscribeFilter struct {
	Page              string
	Limit             string
	Token             string
	InscriptionId     string
	InscriptionNumber string
	State             string
}
