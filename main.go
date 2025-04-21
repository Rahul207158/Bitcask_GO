package main

import (
	"fmt"
	"net/http"

	"github.com/Rahul207158/Bitcask_GO/api"
	"github.com/Rahul207158/Bitcask_GO/kvstore"
)

func main(){
	kvstore.KeyDir =make(map[string]int64)

	http.HandleFunc("/put", func(w http.ResponseWriter, r *http.Request) {
		api.PutHandler(w, r) // Now we can access KeyDir directly in the handler
	})
	http.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		api.GetHandler(w, r) // Now we can access KeyDir directly in the handler
	})

	// Start the server
	fmt.Println("SERVER Started on 8080")
	http.ListenAndServe(":8080", nil)

}