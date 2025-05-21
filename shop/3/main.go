package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	inv := inventory{"shoes": 50, "socks": 5}
	http.HandleFunc("/list", inv.list)
	http.HandleFunc("/price", inv.price)
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}

type inventory map[string]int

func (inv inventory) list(w http.ResponseWriter, r *http.Request) {
	for item, price := range inv {
		fmt.Fprintf(w, "%s: %d\n", item, price)
	}
}

func (inv inventory) price(w http.ResponseWriter, r *http.Request) {
	item := r.URL.Query().Get("item")
	price, ok := inv[item]
	if !ok {
		msg := fmt.Sprintf("no such item: %q", item)
		http.Error(w, msg, http.StatusNotFound) // 404
		return
	}
	fmt.Fprintf(w, "%d\n", price)
}
