package main

import (
	"log"
	"net/http"

	"github.com/abolfazlnorozian/wireguard_exporter_go/internal/metrics"
	"github.com/abolfazlnorozian/wireguard_exporter_go/pkg/config"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	cfg := config.ParseFlags("1.0.0")

	// Register metrics with configuration
	metrics.Register(cfg.Aliases, cfg.Verbose, cfg.Interval)

	http.Handle(cfg.MetricsPath, promhttp.Handler())

	log.Printf("🚀 WireGuard Exporter running on %s%s", cfg.ListenAddress, cfg.MetricsPath)
	log.Printf("⏱️  Collection interval: %v", cfg.Interval)
	if len(cfg.Aliases) > 0 {
		log.Printf("🏷️  Loaded %d peer aliases", len(cfg.Aliases))
	}

	if err := http.ListenAndServe(cfg.ListenAddress, nil); err != nil {
		log.Fatalf("❌ cannot start server: %v", err)
	}
}
