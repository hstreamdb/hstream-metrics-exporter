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

	StreamStats := GetStreamStats()
	for _, stats := range StreamStats {
		vecRef := newGaugeVec(Stream, stats)
		stats := stats
		go func() {
			doCollectFor(vecRef, Stream, stats, *scrapeInterval)
		}()
	}

	SubscriptionStats := GetSubscriptionStats()
	for _, stats := range SubscriptionStats {
		vecRef := newGaugeVec(Subscription, stats)
		stats := stats
		go func() {
			doCollectFor(vecRef, Subscription, stats, *scrapeInterval)
		}()
	}

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

func newGaugeVec(category string, stats Stats) *prometheus.GaugeVec {
	name := category + "_" + stats.methods[0]
	retVec := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: name,
		},
		[]string{name, "type"},
	)
	prometheus.MustRegister(retVec)
	return retVec
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
			for _, interval := range intervals {
				err := doGetSet(category, interval, method, vec)
				if err != nil {
					//log.Printf("error: %v\n", err)
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
		for k, v := range x {
			if k != mainKey {
				v, err := strconv.ParseFloat(v, 64)
				if err != nil {
					return err
				}
				vec.WithLabelValues(mainKey, k).Set(v)
			}
		}
	}
	return nil
}
