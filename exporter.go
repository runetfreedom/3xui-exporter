package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"strconv"
)

type Collector struct {
	inboundReceive  *prometheus.Desc
	inboundTransmit *prometheus.Desc
	userReceive     *prometheus.Desc
	userTransmit    *prometheus.Desc
}

func NewCollector() *Collector {
	return &Collector{
		inboundReceive: prometheus.NewDesc("xui_inbound_receive_bytes_total",
			"Total receive bytes of inbound.",
			[]string{"id"},
			nil,
		),

		inboundTransmit: prometheus.NewDesc("xui_inbound_transmit_bytes_total",
			"Total transmit bytes of inbound.",
			[]string{"id"},
			nil,
		),

		userReceive: prometheus.NewDesc("xui_user_receive_bytes_total",
			"Total receive bytes of user.",
			[]string{"inbound", "email"},
			nil,
		),

		userTransmit: prometheus.NewDesc("xui_user_transmit_bytes_total",
			"Total transmit bytes of user.",
			[]string{"inbound", "email"},
			nil,
		),
	}
}

func (collector *Collector) Describe(ch chan<- *prometheus.Desc) {
	// Add one of these lines for each of your collectors declared above
	ch <- collector.inboundReceive
	ch <- collector.inboundTransmit
	ch <- collector.userReceive
	ch <- collector.userTransmit
}

func (collector *Collector) Collect(ch chan<- prometheus.Metric) {
	inbounds, err := getInbounds()
	if err != nil {
		log.Println("get inbounds error:", err)

		return
	}

	for _, inbound := range inbounds {
		inboundId := strconv.FormatInt(inbound.Id, 10)

		ch <- prometheus.MustNewConstMetric(collector.inboundReceive, prometheus.CounterValue, float64(inbound.Down), inboundId)
		ch <- prometheus.MustNewConstMetric(collector.inboundTransmit, prometheus.CounterValue, float64(inbound.Up), inboundId)

		for _, clientStat := range inbound.ClientStats {
			ch <- prometheus.MustNewConstMetric(collector.userReceive, prometheus.CounterValue, float64(clientStat.Down), inboundId, clientStat.Email)
			ch <- prometheus.MustNewConstMetric(collector.userTransmit, prometheus.CounterValue, float64(clientStat.Up), inboundId, clientStat.Email)
		}
	}
}
