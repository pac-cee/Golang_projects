package main

import (
	"io"
	"log"
	"net/http"
)

func handleProxy(w http.ResponseWriter, r *http.Request) {
	url := r.RequestURI[1:]
	if url == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Usage: /http://example.com or /https://example.com"))
		return
	}
	log.Printf("Proxying request to: %s", url)
	resp, err := http.DefaultClient.Do(&http.Request{
		Method: r.Method,
		URL:    r.URL,
		Header: r.Header,
		Body:   r.Body,
	})
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte("Error forwarding request: " + err.Error()))
		return
	}
	defer resp.Body.Close()
	for k, v := range resp.Header {
		for _, vv := range v {
			w.Header().Add(k, vv)
		}
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func main() {
	log.Println("Starting HTTP proxy on :8080")
	http.HandleFunc("/", handleProxy)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
