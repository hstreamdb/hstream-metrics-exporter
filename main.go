package main

import (
	"flag"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"strconv"
	"strings"
)

const resourceSuffix = "/v1/cluster"

func main() {
	initExporterMetrics()

	const LOCALHOST = "localhost"
	const HttpPrefix = "http://"
	var (
		serverHost = flag.String("host", LOCALHOST, "")
		serverPort = flag.String("port", "9270", "")

		httpServerHost = flag.String("http-server-host", LOCALHOST, "")
		httpServerPort = flag.String("http-server-port", "9290", "")

		scrapeInterval = flag.Int("scrape-interval", 1, "")
	)

	flag.Parse()
	if !strings.HasPrefix(*httpServerHost, HttpPrefix) {
		*httpServerHost = HttpPrefix + *httpServerHost
	}

	serverHostPort := *serverHost + ":" + *serverPort
	httpServerHostPort := *httpServerHost + ":" + *httpServerPort
	resourceUrl := httpServerHostPort + resourceSuffix
	requestBuilder := NewRequestBuilder(resourceUrl)

	doCollectForCategory(requestBuilder, Stream, GetStreamStats(), StreamName, *scrapeInterval)
	doCollectForCategory(requestBuilder, Subscription, GetSubscriptionStats(), SubscriptionId, *scrapeInterval)
	doCollectForCategory(requestBuilder, ServerHistogram, GetServerHistogramStats(), ServerHistogram, *scrapeInterval)
	doCollectForCategory(requestBuilder, StreamCounter, GetStreamCounterStats(), StreamName, *scrapeInterval)

	http.Handle("/metrics", promhttp.HandlerFor(
		prometheus.DefaultGatherer,
		promhttp.HandlerOpts{
			EnableOpenMetrics: true,
		},
	))
	log.Println(serverHostPort + "/metrics")
	log.Fatalln(
		http.ListenAndServe(serverHostPort, nil),
	)
}

func doCollectForCategory(rb RequestBuilder, category string, categoryStats []Stats, mainKey string, scrapeInterval int) {
	for _, stats := range categoryStats {
		vecRef := newGaugeVec(category, stats, mainKey)
		stats := stats
		go func() {
			doCollectFor(rb, vecRef, category, stats, scrapeInterval)
		}()
	}
}

func newGaugeVec(category string, stats Stats, mainKey string) *prometheus.GaugeVec {
	name := category + "_" + stats.methods[0]
	retVec := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: name,
		},
		[]string{mainKey, "type"},
	)
	prometheus.MustRegister(retVec)
	return retVec
}

func doCollectFor(requestBuilder RequestBuilder, vec *prometheus.GaugeVec, category string, stats Stats, scrapeInterval int) {
	ticker := NewTickerSec(scrapeInterval)
	method := stats.methods[0]
	defer ticker.Stop()
	defer fmt.Printf("stop do collect for: %v %v\n", category, method)
	for {
		select {
		case <-ticker.C:

			err := doGetSet(requestBuilder, category, method, vec)
			if err != nil {
				log.Printf("error: %v\n", err)
				log.Println(category, method)
				log.Println()
			}
		}
	}
}

func doGetSet(requestBuilder RequestBuilder, category, metrics string, vec *prometheus.GaugeVec) error {
	url := requestBuilder(category, metrics)
	xs, err := GetVal(url)
	if err != nil {
		return err
	}

	var mainKey string
	if category == Stream || category == StreamCounter {
		mainKey = StreamName
	} else if category == Subscription {
		mainKey = SubscriptionId
	} else {

		for _, x := range xs {
			for k, v := range x {
				v, err := strconv.ParseFloat(v, 64)
				if err != nil {
					return err
				}
				vec.WithLabelValues(category, k).Set(v)
			}
		}

		return nil
	}

	for _, x := range xs {
		mainVal := x[mainKey]
		for k, v := range x {
			if k != mainKey {
				v, err := strconv.ParseFloat(v, 64)
				if err != nil {
					return err
				}
				vec.WithLabelValues(mainVal, k).Set(v)
			}
		}
	}
	return nil
}

func initExporterMetrics() {
	TotalRequestsCnt = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "hstream_metrics_exporter_total_requests",
		},
	)
	prometheus.MustRegister(TotalRequestsCnt)

	FailedRequestsCnt = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "hstream_metrics_exporter_failed_requests",
		},
	)
	prometheus.MustRegister(FailedRequestsCnt)
}
