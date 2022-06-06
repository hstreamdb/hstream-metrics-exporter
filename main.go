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

var requestBuilder RequestBuilder

func main() {
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
	requestBuilder = NewRequestBuilder(resourceUrl)

	doCollectForCategory(Stream, GetStreamStats(), StreamName, *scrapeInterval)
	doCollectForCategory(Subscription, GetSubscriptionStats(), SubscriptionId, *scrapeInterval)

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

func doCollectForCategory(category string, categoryStats []Stats, mainKey string, scrapeInterval int) {
	for _, stats := range categoryStats {
		vecRef := newGaugeVec(category, stats, mainKey)
		stats := stats
		go func() {
			doCollectFor(vecRef, category, stats, scrapeInterval)
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

func checkZeroIntervalStats(stats Stats) bool {
	for _, x := range []Stats{GetStreamAppendInRecords()} {
		for _, xMethod := range x.methods {
			for _, statsMethod := range stats.methods {
				if xMethod == statsMethod {
					return true
				}
			}
		}
	}
	return false
}

func doCollectFor(vec *prometheus.GaugeVec, category string, stats Stats, scrapeInterval int) {
	ticker := NewTickerSec(scrapeInterval)
	method := stats.methods[0]
	defer ticker.Stop()
	defer fmt.Printf("stop do collect for: %v %v\n", category, method)
	for {
		select {
		case <-ticker.C:
			intervals := []string{DefaultIntervalStr}
			if checkZeroIntervalStats(stats) {
				intervals = []string{"0s"}
			}

			for _, interval := range intervals {
				err := doGetSet(category, interval, method, vec)
				if err != nil {
					log.Printf("error: %v\n", err)
					log.Println(category, interval, method)
					log.Println()
				}
			}
		}
	}
}

func doGetSet(category, interval, metrics string, vec *prometheus.GaugeVec) error {
	url := requestBuilder(category, interval, metrics)
	xs, err := GetVal(url)
	if err != nil {
		return err
	}

	var mainKey string
	if category == Stream {
		mainKey = StreamName
	} else if category == Subscription {
		mainKey = SubscriptionId
	} else {
		panic(category)
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
