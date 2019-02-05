package main

import (
	"fmt"
	"net/http"
)

func echo(w http.ResponseWriter, r *http.Request) {
	data := r.URL.Query().Get("key")
	fmt.Printf("%s\n", data)
	fmt.Fprintf(w, "%s", data)
}

func main() {
	http.HandleFunc("/echo", echo)
	http.ListenAndServe(":8080", nil)
}
