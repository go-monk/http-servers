package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	inv := inventory{"shoes": 10, "socks": 5}
	log.Fatal(http.ListenAndServe(":8080", inv))
}

type inventory map[string]int

func (inv inventory) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/list":
		for item, price := range inv {
			fmt.Fprintf(w, "%s: %d\n", item, price)
		}
	case "/price":
		item := r.URL.Query().Get("item")
		price, ok := inv[item]
		if !ok {
			msg := fmt.Sprintf("no such item: %q", item)
			http.Error(w, msg, http.StatusNotFound) // 404
			return
		}
		fmt.Fprintf(w, "%d\n", price)
	}
}
