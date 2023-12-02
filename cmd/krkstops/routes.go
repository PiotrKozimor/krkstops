package main

import (
	_ "embed"
	"log"
	"net/http"
)

//go:embed privacy.txt
var privacy []byte

//go:embed delete_account.txt
var deleteAccount []byte

func Routes() {
	http.HandleFunc("/privacy", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write(privacy)
		if err != nil {
			log.Print(err)
		}
	})
	http.HandleFunc("/delete_account", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write(deleteAccount)
		if err != nil {
			log.Print(err)
		}
	})
}
