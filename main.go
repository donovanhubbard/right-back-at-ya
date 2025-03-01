package main

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const DEFAULT_MESSAGE = "Crab your dog after you pet"

func getStatusCode(url string) (int, error) {
	trimmedUrl := strings.TrimPrefix(url, "/")
	status, err := strconv.Atoi(trimmedUrl)

	if err != nil {
		return 0, errors.New("Not a valid integer code")
	}

	if status >= 100 && status <= 600 {
		return status, nil
	}

	return 0, errors.New("Not a valid status code")
}

func rbay(w http.ResponseWriter, r *http.Request) {
	var returnCode int
	status, error := getStatusCode(r.URL.String())

	if error != nil {
		returnCode = 200
	} else {
		returnCode = status
	}

	fmt.Fprintf(w, "Method: %s\n", r.Method)
	fmt.Fprintf(w, "URL: %v\n", r.URL)
	fmt.Fprintf(w, "Protocol: %s\n", r.Proto)
	fmt.Fprintf(w, "Host: %s\n", r.Host)
	fmt.Fprintf(w, "Remote Address: %s\n", r.RemoteAddr)

	fmt.Fprintf(w, "Headers:\n")

	for key := range r.Header {
		fmt.Fprintf(w, "  %s: %s\n", key, r.Header[key])
	}

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		fmt.Println("Not able to read body. ", err)
	}
	stringBody := string(body)
	if len(stringBody) > 0 {
		fmt.Fprintf(w, "Body:\n")
		fmt.Fprintf(w, "  %v\n", string(body))
	}

	if len(r.Cookies()) > 0 {
		fmt.Fprintf(w, "Cookies:\n")

		for _, c := range r.Cookies() {
			fmt.Fprintf(w, "  Name: %s\n", c.Name)
			fmt.Fprintf(w, "  Value: %s\n", c.Value)
			fmt.Fprintf(w, "  Path: %s\n\n", c.Path)
		}
	}

	fmt.Fprintf(w, "HTTP Status code: %d\n", returnCode)

	slog.Info(fmt.Sprintf("%s %d %s %s", r.RemoteAddr, returnCode, r.Method, r.URL))
}

func message(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	message, messageSet := os.LookupEnv("MESSAGE")
	if !messageSet {
		message = DEFAULT_MESSAGE
	}
	fmt.Fprintf(w, "%s", message)
	slog.Info(fmt.Sprintf("%s %d %s %s", r.RemoteAddr, 200, r.Method, r.URL))
}

func main() {
	port, portSet := os.LookupEnv("PORT")
	if !portSet {
		port = "8080"
	}

	fmt.Printf("Server started on port %s\n", port)
	http.HandleFunc("/", rbay)
	http.HandleFunc("/message", message)
	err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%s", port), nil)
	slog.Error(fmt.Sprintf("%v", err))
}
