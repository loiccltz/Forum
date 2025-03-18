package backend

import (
	"log"
	"net/http"
)

func StartSecureServer(handler http.Handler) {
	server := &http.Server{
		Addr:    ":443",
		Handler: handler,
	}

	log.Fatal(server.ListenAndServeTLS("localhost+2.pem", "localhost+2-key.pem"))
}
