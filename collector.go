package main

import (
	"github.com/prometheus/client_golang/prometheus"
)

var totalReq = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name: "http_request_total",
	Help: "Total number of  requests by HTTP code.",
}, []string{"code", "url", "method"})

func init() {

	prometheus.MustRegister(totalReq)
}
