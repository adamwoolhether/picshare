package main

import (
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html") //you can test "text/plain" to see result
	//fmt.Fprint(w, r.URL.Path)
	if r.URL.Path == "/" {
		fmt.Fprint(w, "<h1>Welcome the the picapp site</h1>")
	} else if r.URL.Path == "/contact" {
	fmt.Fprint(w, "<a>Contact me: <a href=\"mailto:adamwoolhether@gmail.com\">adamwoolhether@gmail.com</a></h1>")
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "<h1>This page doesn't exist :o</h1><p>Email me if you keep getting this message.</p>")
	}
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe("localhost:3000", nil)

	// Under the hood, this uses ServeMux router:
	/*	mgit ux := &http.ServeMux{}
		mux.HandleFunc("/", handler)
		http.ListenAndServe("localhost:3000", mux)*/
}
