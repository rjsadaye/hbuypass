package main

import (
	"fmt"
	"github.com/hashgraph/hedera-sdk-go"
	"time"
)

func main() {
	// Target account to get the balance for
	accountID := hedera.AccountID{Account:1002}

	client, err := hedera.Dial("testnet.hedera.com:51003")
	if err != nil {
		panic(err)
	}

	client.SetNode(hedera.AccountID{Account: 3})
	client.SetOperator(accountID, func() hedera.SecretKey {
		operatorSecret, err := hedera.SecretKeyFromString("302e020100300506032b6570042204208e8279d0f33ad7ee1a7f7a1ddb6da7d2147c014ca07f54d517b53a3d6804b3a7")
		if err != nil {
			panic(err)
		}

		return operatorSecret
	})

	defer client.Close()

	// Get the _answer_ for the query of getting the account balance
	balance, err := client.Account(accountID).Balance().Get()
	if err != nil {
		panic(err)
	}

	fmt.Printf("balance = %v tinybars\n", balance)
	fmt.Printf("balance = %.5f hbars\n", float64(balance)/100000000.0)
	time.Sleep(1 * time.Second)
}


