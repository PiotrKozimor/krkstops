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

//go:embed legal_notice.txt
var legalNotice []byte

func Routes() {
	handler := func(data []byte) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Cache-Control", "max-age=86400")
			_, err := w.Write(data)
			if err != nil {
				log.Print(err)
			}
		}
	}
	http.HandleFunc("/privacy", handler(privacy))
	http.HandleFunc("/delete_account", handler(deleteAccount))
	http.HandleFunc("/legal_notice", handler(legalNotice))
}
