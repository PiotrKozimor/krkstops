package main

import (
	_ "embed"
	"log"
	"net/http"
)

//go:embed privacy.txt
var privacy []byte

func HandlePrivacy() {
	http.HandleFunc("/privacy", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write(privacy)
		if err != nil {
			log.Print(err)
		}
	})

}
