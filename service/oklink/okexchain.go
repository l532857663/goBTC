package oklink

import "sync"

type Platform struct {
	HttpHeader map[string]string
	// Client     okexChainAPI
}

const (
	OkLinkApiHost         = "https://www.oklink.com"
	GetBalanceForSymbol   = "/api/v5/explorer/address/address-summary"
	GetBalanceForToken    = "/api/v5/explorer/address/address-balance-fills"
	GetAddressTransaction = "/api/v5/explorer/address/transaction-list"
	GetTransactionForUTXO = "/api/v5/explorer/address/unspent"
	GetBlockTransaction   = "/api/v5/explorer/block/transaction-list"
	GetInscriptionList    = "/api/v5/explorer/btc/inscriptions-list"
)

func NewPlatformInfo() *Platform {
	httpHeader := map[string]string{
		"Ok-Access-Key": "21140049-abba-4a24-9550-3a37fe4a69c6",
	}
	p := &Platform{
		HttpHeader: httpHeader,
	}
	return p
}

func (p *Platform) Info() string {
	return "OkLink"
}

func (p *Platform) Close(wg *sync.WaitGroup) {
}
