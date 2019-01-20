package main

import (
	"fmt"
	"github.com/hashgraph/hedera-sdk-go"
	//"os"
	"time"
	//"html/template"
	"log"
	//"io"
	"net/http"
	"github.com/GeertJohan/go.rice"
    "github.com/gorilla/mux"
)
//var tmpl = template.Must(template.New("tmpl").ParseFiles("radio.html"))

func dummyhandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w,"You called me")
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w,"You called me")
}
func main() {
	router := mux.NewRouter()
	router.PathPrefix("/").Handler(http.FileServer(rice.MustFindBox("web").HTTPBox()))
	//router.HandleFunc("/callme",HomeHandler).Methods("GET").Schemes("http")
	//router.HandleFunc("/callme", dummyhandler)
    //log.Fatal(http.ListenAndServe(":9000", router))
	
	router.HandleFunc("/callme", func(w http.ResponseWriter, r *http.Request) {

		fmt.Fprintf(w, "You've requested")
	})
	
	srv := &http.Server{
        Handler:      router,
        Addr:         "localhost:9000",
        // Good practice: enforce timeouts for servers you create!
        WriteTimeout: 15 * time.Second,
        ReadTimeout:  15 * time.Second,
    }

    log.Fatal(srv.ListenAndServe())
	
	



	/*
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if err := tmpl.ExecuteTemplate(w, "radio.html", nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	*/

//	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
//	http.Handle("/static/", http.FileServer(http.Dir("fonts/")))
//	http.Handle("/static/", http.FileServer(http.Dir("images/")))
//	http.Handle("/static/", http.FileServer(http.Dir("js/")))
//	http.Handle("/static/", http.FileServer(http.Dir("media/")))
//	http.Handle("/static/", http.FileServer(http.Dir("video/")))








	// Serve /callme with a text response.
	
	

	// Start the server at http://localhost:9000
	//log.Fatal(http.ListenAndServe(":9000",http.FileServer(http.Dir("/web"))))



	//js.Global.Set("main", main)
	// Read and decode the operator secret key
	operatorAccountID := hedera.AccountID{Account: 1002}
	operatorSecret, err := hedera.SecretKeyFromString("302e020100300506032b6570042204208e8279d0f33ad7ee1a7f7a1ddb6da7d2147c014ca07f54d517b53a3d6804b3a7")
	if err != nil {
		panic(err)
	}

	// Read and decode target account
	targetAccountID, err := hedera.AccountIDFromString("0:0:1005")
	if err != nil {
		panic(err)
	}

	targetAccountID1,err:= hedera.AccountIDFromString("0:0:1003")
	if err!=nil {

		panic(err)
	}

	//
	// Connect to Hedera
	//

	client, err := hedera.Dial("testnet.hedera.com:51003")
	if err != nil {
		panic(err)
	}

	client.SetNode(hedera.AccountID{Account: 3})
	client.SetOperator(operatorAccountID, func() hedera.SecretKey {
		return operatorSecret
	})

	defer client.Close()

	//
	// Get balance for target account
	//

	balance, err := client.Account(targetAccountID).Balance().Get()
	if err != nil {
		panic(err)
	}

	fmt.Printf("account balance = %v\n", balance)

	//
	// Transfer 100 cryptos to target
	//

	nodeAccountID := hedera.AccountID{Account: 3}
	response, err := client.TransferCrypto().
		// Move 100 out of operator account
		Transfer(operatorAccountID, -500).
		// And place in our new account
		Transfer(targetAccountID, 10).
		Transfer(targetAccountID1,490).
		Operator(operatorAccountID).
		Node(nodeAccountID).
		Memo("[test] hedera-sdk-go v2").
		Sign(operatorSecret). // Sign it once as operator
		Sign(operatorSecret). // And again as sender
		Execute()

	if err != nil {
		panic(err)
	}

	transactionID := response.ID
	fmt.Printf("transferred; transaction = %v\n", transactionID)

	//
	// Get receipt to prove we sent ok
	//

	fmt.Printf("wait for 2s...\n")
	time.Sleep(2 * time.Second)

	receipt, err := client.Transaction(*transactionID).Receipt().Get()
	if err != nil {
		panic(err)
	}

	if receipt.Status != hedera.StatusSuccess {
		panic(fmt.Errorf("transaction has a non-successful status: %v", receipt.Status.String()))
	}

	fmt.Printf("wait for 2s...\n")
	time.Sleep(2 * time.Second)

	//
	// Get balance for target account (again)
	//
	balance, err = client.Account(targetAccountID1).Balance().Get()
	if err != nil {
		panic(err)
	}

	fmt.Printf("account1 balance = %v\n", balance)

	balance, err = client.Account(targetAccountID).Balance().Get()
	if err != nil {
		panic(err)
	}

	fmt.Printf("account2 balance = %v\n", balance)
}