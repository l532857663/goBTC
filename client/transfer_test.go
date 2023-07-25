package client

import (
	"fmt"
	"testing"
)

func Test_GetTransactionInfo(t *testing.T) {
	transferHash := "7fb631b7ed420c07b546ee4db8527a9523bbc44961f9983430166988cd6beeeb"
	txInfo, err := srv.GetTransactionByHash(transferHash)
	if err != nil {
		fmt.Printf("GetTransactionInfo err: %+v, txHash: %+v\n", err, transferHash)
		return
	}
	fmt.Printf("txInfo: %+v\n", txInfo)
	fmt.Printf("txInfo msg: %+v\n", txInfo.MsgTx())
}

func Test_GetRawTransactionInfo(t *testing.T) {
	transferHash := "7fb631b7ed420c07b546ee4db8527a9523bbc44961f9983430166988cd6beeeb"
	txInfo, err := srv.GetRawTransactionByHash(transferHash)
	if err != nil {
		fmt.Printf("GetTransactionInfo err: %+v, txHash: %+v\n", err, transferHash)
		return
	}
	fmt.Printf("txInfo: %+v\n", txInfo)
}
