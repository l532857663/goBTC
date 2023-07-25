package service

import (
	"goBTC/models"
	"sync"
)

type (
	Platform interface {
		Info() string
		Close(wg *sync.WaitGroup)
	}

	GetBalance interface {
		Platform
		GetBalanceByAddress(symbol, address, protocolType string, filter models.Filter) (*models.GetBalanceResp, error)
	}

	GetTransfer interface {
		Platform
		GetTransferByAddress(symbol, address, protocolType string, filter models.Filter) (*models.GetTransferResp, error)
		GetTransferByBlockNum(symbol, height, protocolType string, filter models.Filter) (*models.GetTransferResp, error)
	}

	GetTransferUTXO interface {
		Platform
		GetTransferUTXOByAddress(symbol, address string) ([]*models.GetTransferUTXOResp, error)
	}

	GetInscriptions interface {
		Platform
		GetInscriptionList(filter models.InscribeFilter) (*models.GetInscribeResp, error)
	}

	Platforms map[string]Platform
)
