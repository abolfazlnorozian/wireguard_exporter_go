package metrics

import (
	"log"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"golang.zx2c4.com/wireguard/wgctrl"
)

var (
	// Global metrics
	wireguardInterfacesTotal = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "wireguard_interfaces_total",
		Help: "Total number of interfaces",
	})

	wireguardScrapeSuccess = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "wireguard_scrape_success",
		Help: "If the scrape was a success",
	})

	wireguardScrapeDuration = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "wireguard_scrape_duration_milliseconds",
		Help: "Duration in milliseconds of the scrape",
	})

	// Interface-level metrics
	wireguardPeersTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "wireguard_peers_total",
			Help: "Total number of peers per interfaces",
		},
		[]string{"interface"},
	)

	wireguardBytesTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "wireguard_bytes_total",
			Help: "Total number of bytes per direction per interface",
		},
		[]string{"interface", "direction"},
	)

	// Peer-level metrics
	wireguardPeerBytesTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "wireguard_peer_bytes_total",
			Help: "Total number of bytes per direction for a peer",
		},
		[]string{"interface", "peer", "alias", "direction"},
	)

	wireguardDurationSinceLatestHandshake = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "wireguard_duration_since_latest_handshake",
			Help: "Duration since latest handshake for a peer",
		},
		[]string{"interface", "peer", "alias"},
	)

	wireguardPeerEndpoint = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "wireguard_peer_endpoint",
			Help: "Peers info. static value",
		},
		[]string{"interface", "peer", "alias", "endpoint_ip"},
	)
)

var (
	aliases  map[string]string
	verbose  bool
	interval time.Duration
)

func Register(aliasMap map[string]string, verboseLogging bool, scrapeInterval time.Duration) {
	aliases = aliasMap
	verbose = verboseLogging
	interval = scrapeInterval

	prometheus.MustRegister(
		wireguardInterfacesTotal,
		wireguardScrapeSuccess,
		wireguardScrapeDuration,
		wireguardPeersTotal,
		wireguardBytesTotal,
		wireguardPeerBytesTotal,
		wireguardDurationSinceLatestHandshake,
		wireguardPeerEndpoint,
	)

	go collect()
}

func collect() {
	client, err := wgctrl.New()
	if err != nil {
		log.Printf("❌ cannot open wgctrl client: %v", err)
		wireguardScrapeSuccess.Set(0)
		return
	}
	defer client.Close()

	// ticker for updating regularly
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		collectOnce(client)
		<-ticker.C
	}
}

func collectOnce(client *wgctrl.Client) {
	start := time.Now()

	devices, err := client.Devices()
	if err != nil {
		log.Printf("❌ cannot list wireguard devices: %v", err)
		wireguardScrapeSuccess.Set(0)
		return
	}

	wireguardScrapeSuccess.Set(1)
	wireguardInterfacesTotal.Set(float64(len(devices)))

	for _, dev := range devices {
		if verbose {
			log.Printf("📡 Collecting metrics for interface: %s", dev.Name)
		}

		// Count peers
		wireguardPeersTotal.WithLabelValues(dev.Name).Set(float64(len(dev.Peers)))

		// Calculate total bytes for the interface
		var totalRx, totalTx int64
		for _, peer := range dev.Peers {
			totalRx += peer.ReceiveBytes
			totalTx += peer.TransmitBytes
		}
		wireguardBytesTotal.WithLabelValues(dev.Name, "rx").Set(float64(totalRx))
		wireguardBytesTotal.WithLabelValues(dev.Name, "tx").Set(float64(totalTx))

		// Per-peer metrics
		for _, peer := range dev.Peers {
			pubKey := peer.PublicKey.String()
			alias := getAlias(pubKey)

			if verbose {
				log.Printf("  👤 Peer: %s (alias: %s)", pubKey, alias)
			}

			// Peer bytes
			wireguardPeerBytesTotal.WithLabelValues(dev.Name, pubKey, alias, "rx").Set(float64(peer.ReceiveBytes))
			wireguardPeerBytesTotal.WithLabelValues(dev.Name, pubKey, alias, "tx").Set(float64(peer.TransmitBytes))

			// Duration since last handshake
			var secondsSinceHandshake float64
			if !peer.LastHandshakeTime.IsZero() {
				secondsSinceHandshake = time.Since(peer.LastHandshakeTime).Seconds()
			}
			wireguardDurationSinceLatestHandshake.WithLabelValues(dev.Name, pubKey, alias).Set(secondsSinceHandshake)

			// Peer endpoint
			if peer.Endpoint != nil {
				endpointIP := peer.Endpoint.IP.String()
				wireguardPeerEndpoint.WithLabelValues(dev.Name, pubKey, alias, endpointIP).Set(1)
			}
		}
	}

	// Record scrape duration
	duration := time.Since(start)
	wireguardScrapeDuration.Set(float64(duration.Milliseconds()))

	if verbose {
		log.Printf("✅ Scrape completed in %v", duration)
	}
}

func getAlias(publicKey string) string {
	if alias, ok := aliases[publicKey]; ok {
		return alias
	}
	return ""
}
