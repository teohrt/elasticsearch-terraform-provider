package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

const port = ":8080"

func main() {
	fmt.Printf("Listening on port %s ...\n", port)
	http.HandleFunc("/", handler)
	if err := http.ListenAndServe(port, nil); err != nil {
		panic(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	url := r.URL.Path
	fmt.Printf("=================\nBody: %s\nURL: %s\n", string(body), url)
}
