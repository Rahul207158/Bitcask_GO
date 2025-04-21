package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Rahul207158/Bitcask_GO/kvstore"
)

type Server struct {
	Store *kvstore.Store
}

func NewServer(store *kvstore.Store) *Server {
	return &Server{Store: store}
}

func (s *Server) PutHandler(w http.ResponseWriter, r *http.Request) {
	var payload kvstore.RequestPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	if payload.Key == "" || payload.Value == "" {
		http.Error(w, "Missing key or value", http.StatusBadRequest)
		return
	}

	err = s.Store.Put(payload.Key, payload.Value)
	if err != nil {
		fmt.Println("Error putting value:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Entry added successfully\n")
}

func (s *Server) GetHandler(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	if key == "" {
		http.Error(w, "Key is required", http.StatusBadRequest)
		return
	}

	value, err := s.Store.Get(key)
	if err != nil {
		if err.Error() == fmt.Sprintf("key not found: %s", key) {
			http.Error(w, "Key not found", http.StatusNotFound)
		} else {
			fmt.Println("Error getting value:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	fmt.Fprintf(w, "Value for key '%s': %s\n", key, value)
}
