package config

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

type Config struct {
	ListenAddress string            //port to listen on
	MetricsPath   string            //path endpoint for metrics
	Interval      time.Duration     //get information from wireguard every interval
	GeoIPPath     string            //path to GeoIP2 or GeoLite2 Country database
	Aliases       map[string]string //publicKey:alias entries
	Verbose       bool
	Version       bool
	VersionText   string // for linker
}

func ParseFlags(version string) *Config {
	cfg := &Config{
		VersionText: version,
	}

	var aliasList string

	flag.StringVar(&cfg.ListenAddress, "listen-address", ":9586", "Address to listen on for metrics server")
	flag.StringVar(&cfg.MetricsPath, "metrics-path", "/metrics", "Path under which to expose metrics")
	flag.DurationVar(&cfg.Interval, "interval", 15*time.Second, "Interval between metric collections")
	flag.StringVar(&cfg.GeoIPPath, "geoip-path", "", "Path to GeoIP2 or GeoLite2 Country database")
	flag.StringVar(&aliasList, "alias", "", "Comma-separated list of publicKey:alias entries")
	flag.BoolVar(&cfg.Verbose, "verbose", false, "Enable verbose logging")
	flag.BoolVar(&cfg.Version, "version", false, "Show version and exit")

	flag.Parse()

	if cfg.Version {
		fmt.Printf("WireGuard Exporter Go version: %s\n", cfg.VersionText)
		os.Exit(0)
	}

	cfg.Aliases = parseAliases(aliasList)
	return cfg
}

func parseAliases(aliasList string) map[string]string {
	aliases := make(map[string]string)
	if aliasList == "" {
		return aliases
	}
	for _, pair := range strings.Split(aliasList, ",") {
		kv := strings.Split(pair, ":")
		if len(kv) == 2 {
			aliases[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
		}
	}
	return aliases
}
