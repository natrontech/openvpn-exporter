package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/natrontech/openvpn-exporter/exporters"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// These variables are set in build step
var Version = "v0.0.0-dev.0"
var Commit = "none"
var BuildTime = "unknown"

func main() {
	var (
		statusfiles = flag.String("openvpn.status-files", "/var/run/openvpn-status.log",
			"Path to OpenVPN status files")
		metricsPath = flag.String("openvpn.metrics-path", "/metrics",
			"Path under which to expose metrics")
		listenAddress = flag.String("openvpn.listen-address", ":9176",
			"Address on which to expose metrics")
		loglevel = flag.String("openvpn.loglevel", "info",
			"Loglevel")
		ignoreIndividuals = flag.String("openvpn.ignore-individuals", "false",
			"If ignoring metrics for individuals")
	)
	flag.Parse()

	// log build information
	log.Printf("INFO: Starting OpenVPN Exporter %s, commit %s, built at %s", Version, Commit, BuildTime)

	// if env variable is set, it will overwrite defaults or flags
	if os.Getenv("OPENVPN_LOGLEVEL") != "" {
		*loglevel = os.Getenv("OPENVPN_LOGLEVEL")
	}
	if os.Getenv("OPENVPN_STATUS_FILES") != "" {
		*statusfiles = os.Getenv("OPENVPN_STATUS_FILES")
	}
	if os.Getenv("OPENVPN_METRICS_PATH") != "" {
		*metricsPath = os.Getenv("OPENVPN_METRICS_PATH")
	}
	if os.Getenv("OPENVPN_LISTEN_ADDRESS") != "" {
		*listenAddress = os.Getenv("OPENVPN_LISTEN_ADDRESS")
	}
	if os.Getenv("OPENVPN_IGNORE_INDIVIDUALS") != "" {
		*ignoreIndividuals = os.Getenv("OPENVPN_IGNORE_INDIVIDUALS")
	}

	// convert flags
	ignoreIndividualsBool, err := strconv.ParseBool(*ignoreIndividuals)
	if err != nil {
		log.Fatalf("ERROR: Unable to parse ignoreIndividuals: %s", err)
	}

	files := strings.Split(*statusfiles, ",")

	// check if statusfile exists
	for _, file := range files {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			log.Fatalf("ERROR: Status file does not exist: %s", file)
		}
	}

	// debug
	if *loglevel == "debug" {
		log.Printf("DEBUG: Using status files: %s", *statusfiles)
		log.Printf("DEBUG: Using metrics path: %s", *metricsPath)
		log.Printf("DEBUG: Using listen address: %s", *listenAddress)
		log.Printf("DEBUG: Ignoring individuals: %t", ignoreIndividualsBool)
	}

	log.Printf("INFO: Listening on: %s", *listenAddress)
	log.Printf("INFO: Metrics path: %s", *metricsPath)

	exporter, err := exporters.NewOpenVPNExporter(files, ignoreIndividualsBool)
	if err != nil {
		panic(err)
	}
	prometheus.MustRegister(exporter)

	http.Handle(*metricsPath, promhttp.Handler())
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
