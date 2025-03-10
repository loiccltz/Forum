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

	log.Fatal(server.ListenAndServeTLS("server.crt", "server.key"))
}
