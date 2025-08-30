[net/http](https://pkg.go.dev/net/http) is the standard's library package for writing HTTP servers (and clients). It has two important components: the  `ListenAndServe` function and the `Handler` interface:

```go
package http

func ListenAndServe(address string, h Handler) error

type Handler interface {
	ServeHTTP(w ResponseWriter, r *Request)
}
```

ListenAndServe runs forever until it fails (always with a non-nil error). It requires an instance of the Handler interface to which all requests will be dispatched (routed).

## shop 1

Let's work on a simple e-shop selling shoes and socks. We need a Handler that will take in the requests and generate responses. It can be any type that implements the ServeHTTP method. In our case it will be a map of strings to integers representing items and their prices:

```go
type inventory map[string]int

func (inv inventory) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for item, price := range inv {
		fmt.Fprintf(w, "%s: %d\n", item, price)
	}
}
```

Now we initialize the inventory with some data and pass it to ListenAndServe:

```go
func main() {
	inv := inventory{"shoes": 10, "socks": 5}
	log.Fatal(http.ListenAndServe(":8080", inv))
}
```

Let's run it:

```sh
❯ go run ./shop/1/main.go &
❯ curl localhost:8080
shoes: 10
socks: 5
```

Nice!

## shop 2

Note that the server will list the inventory (i.e. run the ServeHTTP method of the inventory type) for every request, regardless of URL. We might want to get different data based on the path component of the URL. Let's implement a second version of the shop by modifying the ServeHTTP function:

```go
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
```

Does it work?

```sh
❯ go run ./shop/2/main.go &
❯ curl localhost:8080/list
shoes: 10
socks: 5
❯ curl localhost:8080/price?item=socks
5
```

Yep!

## shop 3

Instead of looking at the URL path inside the ServeHTTP we can use the `ServeMux`, a request multiplexer, to simplify the association between URLs and handlers. First let's split the functionality we have in ServeHTTP into two functions:

```go
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
```

Nor list neither price satisfy the Handler interface now since there's no more the ServeHTTP method. But they implement handler-like behavior because they handle requests and provide responses. We just need to *convert* them to a proper Handler by using http.HandlerFunc adapter:

```go
func main() {
	inv := inventory{"shoes": 10, "socks": 5}
	mux := http.NewServeMux()
	mux.Handle("/list", http.HandlerFunc(inv.list))
	mux.Handle("/price", http.HandlerFunc(inv.price))
	log.Fatal(http.ListenAndServe(":8080", mux))
}
```

Because registering a handler this way is so common, ServeMux has a convenience method HandleFunc:

```go
func main() {
	db := database{"shoes": 50, "socks": 5}
	mux := http.NewServeMux()
	mux.HandleFunc("/list", db.list)
	mux.HandleFunc("/price", db.price)
	log.Fatal(http.ListenAndServe("localhost:8000", mux))
}
```

Also for convenience net/http provides a global ServeMux instance called DefaultServeMux. To use DefaultServeMux pass nil to ListenAndServe:

```go
func main() {
	db := database{"shoes": 50, "socks": 5}
	http.HandleFunc("/list", db.list)
	http.HandleFunc("/price", db.price)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}
```
