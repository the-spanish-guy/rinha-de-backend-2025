package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Teste")
	server := http.NewServeMux()
	server.HandleFunc("POST /payments", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("deu bom")
		w.WriteHeader(http.StatusAccepted)
	})

	http.ListenAndServe("localhost:8080", server)

	fmt.Print("alo")

}
