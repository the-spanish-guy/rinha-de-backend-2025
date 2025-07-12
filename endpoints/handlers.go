package endpoints

import (
	"fmt"
	"net/http"
)

func PaymentHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("deu bom")
	w.WriteHeader(http.StatusAccepted)
}

func PaymentSummaryHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("deu bom")
	w.WriteHeader(http.StatusAccepted)
}
