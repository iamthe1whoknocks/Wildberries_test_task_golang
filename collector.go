package main

import (
	"github.com/prometheus/client_golang/prometheus"
)

var totalReq = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name: "http_request_total",
	Help: "Total number of  requests by HTTP code.",
}, []string{"code", "url", "method"})

var reqDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
	Name: "http_request_duration",
	Help: "Duration of requests by HTTP code.",
}, []string{"code", "url", "method"})

func init() {

	prometheus.MustRegister(totalReq)
	prometheus.MustRegister(reqDuration)
}
