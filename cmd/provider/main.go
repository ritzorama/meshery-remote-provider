package main

import (
	"log"
	"net/http"

	"github.com/tata-consulting/meshery-remote-provider/internal/provider"
)

func main() {
	cfg := provider.LoadConfig()

	server := provider.NewServer(cfg)
	addr := ":" + cfg.Port

	log.Printf("meshery remote provider listening on %s", addr)
	if err := http.ListenAndServe(addr, server); err != nil {
		log.Fatal(err)
	}
}
