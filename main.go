package main

import (
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html") //you can test "text/plain" to see result
	fmt.Fprint(w, "<h1>Welcome the the picapp site</h1>")
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe("localhost:3000", nil)
}
