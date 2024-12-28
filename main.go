package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const promNamespace = "openvpn"

// These variables are set in build step
var Version = "v0.0.0-dev.0"
var Commit = "none"
var BuildTime = "unknown"

var (
	// Flags
	statusfile = flag.String("status-file", "/var/run/openvpn-status.log",
		"Path to OpenVPN status file")
	metricsPath = flag.String("openvpn.metrics-path", "/metrics",
		"Path under which to expose metrics")
	listenAddress = flag.String("openvpn.listen-address", ":9999",
		"Address on which to expose metrics")
	loglevel = flag.String("openvpn.loglevel", "info",
		"Loglevel")

	// Metrics
	up = prometheus.NewDesc(
		prometheus.BuildFQName(promNamespace, "", "up"),
		"Was the last query of OpenVPN Exporter successful.",
		nil, nil,
	)
)

type Exporter struct {
	statusfile string
}

func NewExporter(statusfile string) *Exporter {
	return &Exporter{
		statusfile: statusfile,
	}
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- up
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	err := e.collectFromAPI(ch)
	if err != nil {
		ch <- prometheus.MustNewConstMetric(
			up, prometheus.GaugeValue, 0,
		)
		log.Println(err)
		return
	}
	ch <- prometheus.MustNewConstMetric(
		up, prometheus.GaugeValue, 1,
	)
}

func (e *Exporter) collectFromAPI(ch chan<- prometheus.Metric) error {

	return nil
}

func main() {
	flag.Parse()

	// log build information
	log.Printf("INFO: Starting OpenVPN Exporter %s, commit %s, built at %s", Version, Commit, BuildTime)

	// if env variable is set, it will overwrite defaults or flags
	if os.Getenv("OPENVPN_LOGLEVEL") != "" {
		*loglevel = os.Getenv("OPENVPN_LOGLEVEL")
	}
	if os.Getenv("OPENVPN_STATUS_FILE") != "" {
		*statusfile = os.Getenv("OPENVPN_STATUS_FILE")
	}
	if os.Getenv("OPENVPN_METRICS_PATH") != "" {
		*metricsPath = os.Getenv("OPENVPN_METRICS_PATH")
	}
	if os.Getenv("OPENVPN_LISTEN_ADDRESS") != "" {
		*listenAddress = os.Getenv("OPENVPN_LISTEN_ADDRESS")
	}

	// debug
	if *loglevel == "debug" {
		log.Printf("DEBUG: Using status file: %s", *statusfile)
		log.Printf("DEBUG: Using metrics path: %s", *metricsPath)
		log.Printf("DEBUG: Using listen address: %s", *listenAddress)
	}

	log.Printf("INFO: Listening on: %s", *listenAddress)
	log.Printf("INFO: Metrics path: %s", *metricsPath)

	// start http server
	http.HandleFunc(*metricsPath, func(w http.ResponseWriter, r *http.Request) {

		exporter := NewExporter(*statusfile)

		// catch if register of exporter fails
		err := prometheus.Register(exporter)
		if err != nil {
			// if register fails, we log the error and return
			log.Printf("ERROR: %s", err)
		}
		promhttp.Handler().ServeHTTP(w, r) // Serve the metrics
		prometheus.Unregister(exporter)    // Clean up after serving
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(`<html>
			<head><title>OpenVPN Exporter</title></head>
			<body>
			<h1>OpenVPN Exporter</h1>
			<p><a href='` + *metricsPath + `'>Metrics</a></p>
			</body>
			</html>`))
		if err != nil {
			log.Printf("ERROR: Failed to write response: %s", err)
		}
	})

	server := &http.Server{
		Addr:         *listenAddress,
		Handler:      nil,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}
	log.Fatal(server.ListenAndServe())
}
