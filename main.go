package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	addr := fmt.Sprintf("127.0.0.1:%s", port)

	http.Handle("/", http.FileServer(http.Dir('.')))

	log.Printf("Listening on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
