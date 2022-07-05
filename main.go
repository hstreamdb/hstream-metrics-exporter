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

const versionSuffix = "/v1"
const clusterSuffix = "/cluster"
const statsSuffix = "/stats"

var (
	resourceUrl string
)

func main() {
	initExporterMetrics()

	const LOCALHOST = "localhost"
	const HttpPrefix = "http://"

	serverHost := flag.String("host", LOCALHOST, "")
	serverPort := flag.String("port", "9270", "")

	httpServerHost := flag.String("http-server-host", LOCALHOST, "")
	httpServerPort := flag.String("http-server-port", "9290", "")

	scrapeInterval := flag.Int("scrape-interval", 1, "")

	flag.Parse()
	if !strings.HasPrefix(*httpServerHost, HttpPrefix) {
		*httpServerHost = HttpPrefix + *httpServerHost
	}

	serverHostPort := *serverHost + ":" + *serverPort
	httpServerHostPort := *httpServerHost + ":" + *httpServerPort
	resourceUrl = httpServerHostPort + versionSuffix
	requestBuilder := NewRequestBuilder(resourceUrl + clusterSuffix)

	doCollectForCategory(requestBuilder, Stream, GetStreamStats(), StreamName, *scrapeInterval)
	doCollectForCategory(requestBuilder, Subscription, GetSubscriptionStats(), SubscriptionId, *scrapeInterval)
	doCollectForCategory(requestBuilder, ServerHistogram, GetServerHistogramStats(), ServerHistogram, *scrapeInterval)
	doCollectForCategory(requestBuilder, StreamCounter, GetStreamCounterStats(), StreamName, *scrapeInterval)
	doCollectForTargets(GetStatsTargets(), *scrapeInterval)

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
		metrics := stats.methods[0]
		url := rb(category, metrics)
		vecRef := newGaugeVec(category+"_"+metrics, mainKey)
		go func() {
			doCollectFor(vecRef, url, category, metrics, scrapeInterval)
		}()
	}
}

func doCollectForTargets(targets []StatsTarget, scrapeInterval int) {
	for _, target := range targets {
		name := strings.Replace(target.uri, "/", "_", -1)
		name = name[1:]

		url := resourceUrl + statsSuffix + target.uri

		vecRef := newGaugeVec(name, target.mainKey)
		target := target
		go func() {
			doCollectFor(vecRef, url, target.category, name, scrapeInterval)
		}()
	}
}

func doCollectFor(vec *prometheus.GaugeVec, url, category string, method string, scrapeInterval int) {
	ticker := NewTickerSec(scrapeInterval)
	defer ticker.Stop()
	defer fmt.Printf("stop do collect for: %v %v\n", category, method)
	for {
		select {
		case <-ticker.C:

			err := doGetSet(url, category, vec)
			if err != nil {
				log.Printf("error: %v\n", err)
				log.Println(category, method)
				log.Println()
			}
		}
	}
}

func doGetSet(url, category string, vec *prometheus.GaugeVec) error {

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

func newGaugeVec(name, mainKey string) *prometheus.GaugeVec {
	retVec := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: name,
		},
		[]string{mainKey, "type"},
	)
	prometheus.MustRegister(retVec)
	return retVec
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
