package main

import (
	"flag"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	ListenAddr          = flag.String("listen-address", getEnv("EXPORTER_ADDRESS", ":8080"), "The address to listen on for HTTP requests.")
	RedashAPIBaseURL    = flag.String("redash-api-base-url", getEnv("REDASH_API_BASE_URL", ""), "the base url address of the redash api endpoint")
	RedashAPIKey        = flag.String("redash-api-key", getEnv("REDASH_API_KEY", ""), "the api key used for retrieving redash resources")
	RedashProbeQueryID  = flag.Int("redash-probe-query-id", getEnvInt("REDASH_PROBE_QUERY_ID", 281), "the redash query probe id")
	RedashProbeAlertID  = flag.Int("redash-probe-alert-id", getEnvInt("REDASH_PROBE_ALERT_ID", 42), "the redash alert probe id")
	RedashProbeInterval = flag.Int("redash-probe-interval", getEnvInt("REDASH_PROBE_INTERVAL", 10), "the redash schedular probe interval in minutes")
)

func main() {
	RedashCollector := &RedashCollector{
		AlertStatusDesc: prometheus.NewDesc(
			"redash_alert_status",
			"Alert status of the redash scheduler proble.",
			[]string{"status"},
			nil,
		),
		QueryRefreshStatusDesc: prometheus.NewDesc(
			"redash_query_refresh_status",
			"Query refresh status of the redash scheduler proble",
			[]string{},
			nil,
		),
	}

	reg := prometheus.NewPedanticRegistry()
	reg.MustRegister(RedashCollector)

	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	http.ListenAndServe(*ListenAddr, nil)
}
