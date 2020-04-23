package main

import (
	"flag"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
	"net/http"
)

var (
	listenAddress = flag.String("listen-address", ":1234", "The address to expose the metrics in prometheus format.")
	// httpURL       = flag.String("http-url", "https://httpstat.us/200", "The URL that should be scraped.")
)

const httpURL = "https://httpstat.us/200"

func GetUrlStatus(url string) int {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		return 1
	}
	return 0
}

func GetUrlResponseTime(url string) float64 {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		return 1
	}
	return 0
}

type Exporter struct {
	up           *prometheus.Desc
	responseUp   *prometheus.Desc
	responseTime *prometheus.Desc
}

func NewExporter() *Exporter {
	return &Exporter{
		up: prometheus.NewDesc(
			"up",
			"Whether the exporter is up and running.",
			nil,
			nil,
		),
		responseUp: prometheus.NewDesc(
			"sample_external_url_up",
			"Whether the URL is up and running.",
			[]string{"url"},
			nil,
		),
		responseTime: prometheus.NewDesc(
			"sample_external_url_response_ms",
			"Number time it took for the http request to finish (in milliseconds).",
			[]string{"url"},
			nil,
		),
	}
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	defer func() {
		if err := recover(); err != nil {
			ch <- prometheus.MustNewConstMetric(e.up, prometheus.GaugeValue, 0)
			log.Fatal(err)
		}
	}()

	ch <- prometheus.MustNewConstMetric(e.up, prometheus.GaugeValue, 1)
	ch <- prometheus.MustNewConstMetric(e.responseUp, prometheus.GaugeValue, float64(GetUrlStatus(httpURL)), httpURL)
	ch <- prometheus.MustNewConstMetric(e.responseTime, prometheus.GaugeValue, GetUrlResponseTime(httpURL), httpURL)
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- e.up
	ch <- e.responseUp
	ch <- e.responseTime
}

func main() {
	prometheus.MustRegister(NewExporter())
	http.Handle("/metrics", promhttp.Handler())

	log.Infoln("Listening on", *listenAddress)
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}
