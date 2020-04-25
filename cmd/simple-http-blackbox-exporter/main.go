package main

import (
	"flag"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

var (
	listenAddress = flag.String("listen-address", ":1234", "The address to expose the metrics in prometheus format.")
	configFile    = flag.String("config", "config.yaml", "A yaml config file with settings and services to monitor.")
)

type Config struct {
	Urls []string `yaml:",flow"`
}

func (c *Config) GetConf() *Config {

	if *configFile == "" {
		log.Fatal("Please provide yaml file by using -config option")
	}

	file, err := ioutil.ReadFile(*configFile)
	if err != nil {
		log.Fatalf("Error while reading the YAML file: %s\n", err)
	}

	err = yaml.Unmarshal(file, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return c
}

func GetURL(url string) (int, int64) {
	start := time.Now()

	client := http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Get(url)
	if err != nil {
		log.Error(err)
		return 0, 0
	}

	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		defer resp.Body.Close()

		elapsed := time.Since(start).Milliseconds()
		log.Info("client.Get to " + url + " took " + strconv.FormatInt(elapsed, 10) + " milliseconds.\n")
		return 1, elapsed
	}

	log.Warn("client.Get failed to get to " + url + " - status code: " + strconv.Itoa(resp.StatusCode) + ".\n")
	return 0, 0
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
			"The amount of time it took for the http request to finish (in milliseconds).",
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

	var c Config
	c.GetConf()
	for _, url := range c.Urls {
		var responseUp, responseTime = GetURL(url)
		ch <- prometheus.MustNewConstMetric(e.responseUp, prometheus.GaugeValue, float64(responseUp), url)
		ch <- prometheus.MustNewConstMetric(e.responseTime, prometheus.GaugeValue, float64(responseTime), url)
	}
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- e.up
	ch <- e.responseUp
	ch <- e.responseTime
}

func main() {
	flag.Parse()
	prometheus.MustRegister(NewExporter())
	http.Handle("/metrics", promhttp.Handler())

	log.Infoln("Listening on", *listenAddress)
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}
