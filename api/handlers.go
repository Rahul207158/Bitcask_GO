package api

import (
	"fmt"
	"net/http"
	"time"
	"encoding/json"
	"github.com/Rahul207158/Bitcask_GO/kvstore"
)
func PutHandler(w http.ResponseWriter, r *http.Request) {
	var payload kvstore.RequestPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	fmt.Println("Received:", payload.Key, payload.Value)

	if payload.Key == "" || payload.Value == "" {
		http.Error(w, "Missing key or value", http.StatusBadRequest)
		return
	}

	entry := kvstore.Entry{
		TimeStamp:  time.Now().Unix(),
		Key:        payload.Key,
		Value:      payload.Value,
		KeySize:    int32(len(payload.Key)),
		ValueSize:  int32(len(payload.Value)),
	}

	offset, err := kvstore.WriteEntry("data/store_data", entry)
	if err != nil {
		fmt.Println("ERR is there in write  ", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	kvstore.KeyDir[entry.Key] = offset

	fmt.Fprintf(w, "Entry added successfully with offset %d\n", offset)
}


func GetHandler(w http.ResponseWriter,r *http.Request){
	key := r.URL.Query().Get("key")
	if key == "" {
		http.Error(w, "Key is required", http.StatusBadRequest)
		return
	}
	offset,exist:=kvstore.KeyDir[key]
	if !exist{
		http.Error(w, "Key not found", http.StatusNotFound)
		return
	}
	value,err:=kvstore.ReadEntry("data/store_data",offset)
	if err!=nil{
		fmt.Println("err in reading file",err)
	}
	fmt.Fprintf(w, "Value for key '%s': %s\n", key, value)
}