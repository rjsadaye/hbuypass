package main

import (
	"fmt"
	"github.com/hashgraph/hedera-sdk-go"
	"time"
)

func main() {
	// Target account to get the balance for
	accountID := hedera.AccountID{Account:1003}

	client, err := hedera.Dial("testnet.hedera.com:51003")
	if err != nil {
		panic(err)
	}

	client.SetNode(hedera.AccountID{Account: 3})
	client.SetOperator(accountID, func() hedera.SecretKey {
		operatorSecret, err := hedera.SecretKeyFromString("302e020100300506032b65700422042064a33647f036a5889080f8784d491466bd6ecd0bd8ac5d027aa50f54d4498a08")
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


