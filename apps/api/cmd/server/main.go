// Command server menjalankan scaffold API lokal. Wiring domain ada pada M1-02.
package main

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

func main() {
	server := &http.Server{
		Addr:              "127.0.0.1:8080",
		Handler:           chi.NewRouter(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("Scaffold API lokal berjalan di http://%s; endpoint produk belum tersedia", server.Addr)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal("Server lokal tidak dapat dijalankan: ", err)
	}
}
