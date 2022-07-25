package asn

import (
	"errors"
	"strconv"

	"github.com/lwlcom/fastnetmon_exporter/rpc"

	"github.com/lwlcom/fastnetmon_exporter/collector"
	"github.com/prometheus/client_golang/prometheus"
)

const prefix string = "fastnetmon_asn_"

var (
	incomingPacketsDescV4 *prometheus.Desc
	incomingBytesDescV4   *prometheus.Desc
	incomingPacketsDescV6 *prometheus.Desc
	incomingBytesDescV6   *prometheus.Desc
)

func init() {
	l := []string{"target", "asn_number"}
	incomingPacketsDescV4 = prometheus.NewDesc(prefix+"incoming_packets_v4", "Counter for incoming packets per ASN for IPv4", l, nil)
	incomingBytesDescV4 = prometheus.NewDesc(prefix+"incoming_bytes_v4", "Counter for incoming bytes per ASN for IPv4", l, nil)
	incomingPacketsDescV6 = prometheus.NewDesc(prefix+"incoming_packets_v6", "Counter for incoming packets per ASN for IPv6", l, nil)
	incomingBytesDescV6 = prometheus.NewDesc(prefix+"incoming_bytes_v6", "Counter for incoming bytes per ASN for IPv6", l, nil)
}

type asnCollector struct {
}

// NewCollector creates a new collector
func NewCollector() collector.RPCCollector {
	return &asnCollector{}
}

// Describe describes the metrics
func (*asnCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- incomingPacketsDescV4
	ch <- incomingBytesDescV4
	ch <- incomingPacketsDescV6
	ch <- incomingBytesDescV6
}

// Collect collects metrics from FastNetMon API
func (c *asnCollector) Collect(client *rpc.Client, ch chan<- prometheus.Metric, labelValues []string) error {
	var asnV4 = Response{}
	err := client.RunCommandAndParse("/asn_counters_v4", &asnV4)
	if err != nil {
		return err
	}

	if asnV4.Success == false {
		return errors.New(asnV4.ErrorText)
	}

	for _, counter := range asnV4.Values {
		l := append(labelValues, strconv.FormatInt(counter.ASN, 10))

		ch <- prometheus.MustNewConstMetric(incomingPacketsDescV4, prometheus.GaugeValue, float64(counter.IncomingPackets), l...)
		ch <- prometheus.MustNewConstMetric(incomingBytesDescV4, prometheus.GaugeValue, float64(counter.IncomingBytes), l...)
	}

	var asnV6 = Response{}
	err = client.RunCommandAndParse("/asn_counters_v6", &asnV6)
	if err != nil {
		return err
	}

	if asnV6.Success == false {
		return errors.New(asnV6.ErrorText)
	}

	for _, counter := range asnV6.Values {
		l := append(labelValues, strconv.FormatInt(counter.ASN, 10))

		ch <- prometheus.MustNewConstMetric(incomingPacketsDescV6, prometheus.GaugeValue, float64(counter.IncomingPackets), l...)
		ch <- prometheus.MustNewConstMetric(incomingBytesDescV6, prometheus.GaugeValue, float64(counter.IncomingBytes), l...)
	}

	return nil
}
