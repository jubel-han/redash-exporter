package main

import (
	"github.com/prometheus/client_golang/prometheus"
)

type RedashCollector struct {
	AlertStatusDesc        *prometheus.Desc
	QueryRefreshStatusDesc *prometheus.Desc
}

func (c *RedashCollector) CollectAlertStatus() string {
	alert, _ := getAlert(*RedashProbeAlertID)
	return alert.State
}

func (c *RedashCollector) CollectQueryRefreshStatus() float64 {
	query, _ := getQuery(*RedashProbeQueryID)
	result, _ := getQueryResult(query.LastQueryResultID)
	isFresh, _ := isQueryResultFresh(result.RetrievedAt)
	if isFresh {
		return 1
	}
	return 0
}

func (c *RedashCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.AlertStatusDesc
	ch <- c.QueryRefreshStatusDesc
}

func (c *RedashCollector) Collect(ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(
		c.AlertStatusDesc,
		prometheus.GaugeValue,
		float64(1),
		c.CollectAlertStatus(),
	)

	ch <- prometheus.MustNewConstMetric(
		c.QueryRefreshStatusDesc,
		prometheus.GaugeValue,
		c.CollectQueryRefreshStatus(),
	)
}
