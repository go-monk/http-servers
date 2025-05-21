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
	for item, price := range inv {
		fmt.Fprintf(w, "%s: %d\n", item, price)
	}
}
