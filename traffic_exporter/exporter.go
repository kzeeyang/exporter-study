package main

import (
	"crypto/tls"
	"flag"
	"net/http"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	namespace = "Trafficserver"
)

/*
traffic": {
"ccons": 4304,
"cobj": 33013850,
"cpu": 14.47,
"cused": 29.07,
"hit": 280225,
"load": "1.14,1.40,1.90",
"mem": 72.03,
"reqs": 307648,
"rpc": 1.88,
"rx": 806138577,
"time": 1540284300,
"tpr": 80.37,
"tx": 17016624097,
"uptime": 708202,
}
*/

type Traffic_exporter struct {
	URI    string
	mutex  sync.Mutex
	client *http.Client

	up     *prometheus.Desc
	ccons  prometheus.Gauge
	cobj   prometheus.Counter
	cpu    prometheus.Gauge
	cused  prometheus.Gauge
	hit    *prometheus.Desc
	mem    prometheus.Gauge
	reqs   *prometheus.Desc
	rx     *prometheus.Desc
	tx     *prometheus.Desc
	tpr    *prometheus.Desc
	uptime *prometheus.Desc
}

type Traffic_struct struct {
	Traffic struct {
		Ccons  int     `json:"ccons"`
		Cobj   int     `json:"cobj"`
		CPU    float64 `json:"cpu"`
		Cused  float64 `json:"cused"`
		Hit    int     `json:"hit"`
		Load   string  `json:"load"`
		Mem    float64 `json:"mem"`
		Reqs   int     `json:"reqs"`
		RPC    float64 `json:"rpc"`
		Rx     int     `json:"rx"`
		Time   int     `json:"time"`
		Tpr    float64 `json:"tpr"`
		Tx     int     `json:"tx"`
		Uptime int     `json:"uptime"`
	} `json:"traffic"`
}

func NewExporter(uri string, insecure bool) *Traffic_exporter {
	return &Traffic_exporter{
		URI: uri,
		client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: insecure},
			},
		},

		up: prometheus.NewDesc(prometheus.BuildFQName(namespace, "", "up"),
			"Could the traffic server be reached",
			nil,
			nil),
		ccons: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "ccons",
			Help:      "The current connectting requests.",
		}),
		cobj: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "close_objects_total",
			Help:      "Number of traffic server close objects.",
		}),
		cpu: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "cpu",
			Help:      "The current percentage CPU used in system.",
		}),
		cused: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "cused",
			Help:      "The current percentage traffic server used cpu.",
		}),
		hit: prometheus.NewDesc(prometheus.BuildFQName(namespace, "", "hit"),
			"The total hit request in 5 minutes.",
			nil,
			nil),
		mem: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "mem",
			Help:      "The current percentage traffic server used memeory.",
		}),
		reqs: prometheus.NewDesc(prometheus.BuildFQName(namespace, "", "hit"),
			"The total request in 5 minutes.",
			nil,
			nil),
		rx: prometheus.NewDesc(prometheus.BuildFQName(namespace, "", "rx"),
			"The total receive bytes in 5 minutes.",
			nil,
			nil),
		tx: prometheus.NewDesc(prometheus.BuildFQName(namespace, "", "rx"),
			"The total transport bytes in 5 minutes.",
			nil,
			nil),
		tpr: prometheus.NewDesc(prometheus.BuildFQName(namespace, "", "rx"),
			"Unknown.",
			nil,
			nil),
		uptime: prometheus.NewDesc(prometheus.BuildFQName(namespace, "", "rx"),
			"Current uptime in seconds (*)",
			nil,
			nil),
	}
}

func (e *Traffic_exporter) collect(ch chan<- prometheus.Metric) error {

}

func main() {
	var (
		listeningAddress = flag.String("telemetry.address", ":9117", "Address on which to expose metrics.")
		metricsEndpoint  = flag.String("telemetry.endpoint", "/metrics", "Path under which to expose metrics.")
		scrapeURI        = flag.String("scrape_uri", "http://localhost/server-status/?auto", "URI to apache stub status page.")
		insecure         = flag.Bool("insecure", false, "Ignore server certificate if using https.")
		showVersion      = flag.Bool("version", false, "Print version information.")
	)

}
