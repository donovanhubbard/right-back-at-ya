package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func rbay(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Method: %s\n", r.Method)
	fmt.Fprintf(w, "URL: %v\n", r.URL)
	fmt.Fprintf(w, "Protocol: %s\n", r.Proto)
	fmt.Fprintf(w, "Host: %s\n", r.Host)
	fmt.Fprintf(w, "Headers:\n")

	for key := range r.Header {
		fmt.Fprintf(w, "  %s: %s\n", key, r.Header[key])
	}

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		fmt.Println("Not able to read body. ", err)
	}

	fmt.Fprintf(w, "Body:\n")
	fmt.Fprintf(w, "  %v\n", string(body))
	fmt.Fprintf(w, "Remote Address: %s\n", r.RemoteAddr)
	fmt.Fprintf(w, "Cookies:\n")

	// for c := range r.Cookies() {
	// 	fmt.Fprintf(w, "  %v", c)
	// 	// fmt.Fprintf(w, "  Name:%s", c.Name)
	// 	// fmt.Fprintf(w, "  Value:%s", c.Value)
	// 	// fmt.Fprintf(w, "  Path:%s", c.Path)
	// }
}

func main() {
	fmt.Println("Starting")
	mux := &http.ServeMux{}
	mux.HandleFunc("/", rbay)
	http.ListenAndServe(":8085", mux)
	fmt.Println("Exiting")
}
