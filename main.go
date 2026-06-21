package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	namespace = "edgetpu"
)

var (
	labels    = []string{"name"}
	sysfsRoot = "/sys"
)

type EdgeTPUCollector struct {
	sync.Mutex
	numDevices  prometheus.Gauge
	temperature *prometheus.GaugeVec
}

func NewEdgeTPUCollector() *EdgeTPUCollector {
	return &EdgeTPUCollector{
		numDevices: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "num_devices",
				Help:      "Number of EdgeTPU devices",
			},
		),
		temperature: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "temperature_celsius",
				Help:      "EdgeTPU device temperature in Celsius",
			},
			labels,
		),
	}
}

func (c *EdgeTPUCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.numDevices.Desc()
	c.temperature.Describe(ch)
}

func (c *EdgeTPUCollector) Collect(ch chan<- prometheus.Metric) {
	// Only allow one collection at a time
	c.Lock()
	defer c.Unlock()

	c.temperature.Reset()

	devices := FindEdgeTPUDevices()

	c.numDevices.Set(float64(len(devices)))
	ch <- c.numDevices

	for _, device := range devices {
		temp := device.Temperature()
		// Temperature reading is not supported on all devices; skip unknowns
		if temp > 0.0 {
			c.temperature.WithLabelValues(device.name).Set(temp)
		}
	}

	c.temperature.Collect(ch)
}

func envPort() int {
	if val, ok := os.LookupEnv("PORT"); ok {
		if p, err := strconv.Atoi(val); err == nil {
			return p
		}
	}
	return 8080
}

func main() {
	var port int
	portSet := false

	flag.IntVar(&port, "port", envPort(), "Port to listen on (env: PORT)")
	flag.StringVar(&sysfsRoot, "sysfs", "/sys", "Mountpoint of sysfs instance to scan")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "EdgeTPU Prometheus Exporter\n")
		fmt.Fprintf(os.Stderr, "Usage: edgetpu-exporter [flags]\n\n")
		fmt.Fprintf(os.Stderr, "Environment variables:\n")
		fmt.Fprintf(os.Stderr, "  PORT  Port to listen on (default 8080, overridden by -port flag)\n\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	// Detect if -port was explicitly passed
	flag.Visit(func(f *flag.Flag) {
		if f.Name == "port" {
			portSet = true
		}
	})
	_ = portSet

	prometheus.MustRegister(NewEdgeTPUCollector())

	addr := fmt.Sprintf(":%d", port)
	log.Printf("Listening on %s...\n", addr)

	srv := &http.Server{
		Addr:         addr,
		Handler:      promhttp.Handler(),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Fatalf("ListenAndServe error: %v", srv.ListenAndServe())
}
