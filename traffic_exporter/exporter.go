package main

import (
	"fmt"
	"crypto/tls"	
	"net/http"
	"sync"
	"encoding/json"
	"io/ioutil"	
	"flag"
	"os"

	

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/version"
	"github.com/prometheus/common/log"
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
	scrapeFailures prometheus.Counter
	ccons  prometheus.Gauge
	cobj   *prometheus.Desc
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
			scrapeFailures: prometheus.NewCounter(prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "exporter_scrape_failures_total",
				Help:      "Number of errors while scraping traffic server.",
			}),
		ccons: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "ccons",
			Help:      "The current connectting requests.",
		}),
		cobj: prometheus.NewDesc(prometheus.BuildFQName(namespace, "", "cobj"),
			"Could the traffic server be reached",
			nil,
			nil),
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
		reqs: prometheus.NewDesc(prometheus.BuildFQName(namespace, "", "reqs"),
			"The total request in 5 minutes.",
			nil,
			nil),
		rx: prometheus.NewDesc(prometheus.BuildFQName(namespace, "", "rx"),
			"The total receive bytes in 5 minutes.",
			nil,
			nil),
		tx: prometheus.NewDesc(prometheus.BuildFQName(namespace, "", "tx"),
			"The total transport bytes in 5 minutes.",
			nil,
			nil),
		tpr: prometheus.NewDesc(prometheus.BuildFQName(namespace, "", "tpr"),
			"Unknown.",
			nil,
			nil),
		uptime: prometheus.NewDesc(prometheus.BuildFQName(namespace, "", "uptime"),
			"Current uptime in seconds (*)",
			nil,
			nil),
	}
}

func (e *Traffic_exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- e.up
	ch <- e.cobj
	ch <- e.hit
	ch <- e.reqs
	ch <- e.rx
	ch <- e.tx
	ch <- e.tpr
	ch <- e.uptime
	e.scrapeFailures.Describe(ch)
	e.ccons.Describe(ch)
	e.cpu.Describe(ch)
	e.cused.Describe(ch)
	e.mem.Describe(ch)
}

func (e *Traffic_exporter) collect(ch chan<- prometheus.Metric) error {
	resp, err := e.client.Get(e.URI)
	if err != nil {
		ch <- prometheus.MustNewConstMetric(e.up, prometheus.GaugeValue, 0)
		return fmt.Errorf("Error scraping apache: %v", err)
	}
	ch <- prometheus.MustNewConstMetric(e.up, prometheus.GaugeValue, 1)

	data, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode != 200 {
		if err != nil {
			data = []byte(err.Error())
		}
		return fmt.Errorf("Status %s (%d): %s", resp.Status, resp.StatusCode, data)
	}

	var traffic_struct = &Traffic_struct{}
	if  err := json.Unmarshal(data, &traffic_struct); err != nil {
		return fmt.Errorf("Cannot unmarshal body: %s to json",  data)
	}

	//ccons
	e.ccons.Set(float64(traffic_struct.Traffic.Ccons))
	//cobj   prometheus.Counter
	ch <- prometheus.MustNewConstMetric(e.cobj, prometheus.CounterValue, float64(traffic_struct.Traffic.Cobj))
	//cpu    prometheus.Gauge
	e.cpu.Set(float64(traffic_struct.Traffic.CPU))
	//cused  prometheus.Gauge
	e.cused.Set(float64(traffic_struct.Traffic.Cused))
	//hit    *prometheus.Desc
	ch <- prometheus.MustNewConstMetric(e.hit, prometheus.CounterValue, float64(traffic_struct.Traffic.Hit))
	//mem    prometheus.Gauge
	e.mem.Set(float64(traffic_struct.Traffic.Mem))
	//reqs   *prometheus.Desc
	ch <- prometheus.MustNewConstMetric(e.reqs, prometheus.CounterValue, float64(traffic_struct.Traffic.Reqs))
	//rx     *prometheus.Desc
	ch <- prometheus.MustNewConstMetric(e.rx, prometheus.CounterValue, float64(traffic_struct.Traffic.Rx))
	//tx     *prometheus.Desc
	ch <- prometheus.MustNewConstMetric(e.tx, prometheus.CounterValue, float64(traffic_struct.Traffic.Tx))
	//tpr    *prometheus.Desc
	ch <- prometheus.MustNewConstMetric(e.tpr, prometheus.CounterValue, float64(traffic_struct.Traffic.Tpr))
	//uptime *prometheus.Desc
	ch <- prometheus.MustNewConstMetric(e.uptime, prometheus.CounterValue, float64(traffic_struct.Traffic.Uptime))

	return nil
}

func (e *Traffic_exporter) Collect(ch chan<- prometheus.Metric) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	if err := e.collect(ch); err != nil {
		log.Errorf("Error scraping traffic server: %s", err)
		e.scrapeFailures.Inc()
		e.scrapeFailures.Collect(ch)
	}
	return
}

func main() {
	var (
		listeningAddress = flag.String("telemetry.address", ":8110", "Address on which to expose metrics.")
		metricsEndpoint  = flag.String("telemetry.endpoint", "/metrics", "Path under which to expose metrics.")
		scrapeURI        = flag.String("scrape_uri", "http://localhost/_billing", "URI to apache stub status page.")
		insecure         = flag.Bool("insecure", false, "Ignore server certificate if using https.")
		showVersion      = flag.Bool("version", false, "Print version information.")
	)
	
	flag.Parse()

	if *showVersion {
		fmt.Fprintln(os.Stdout, version.Print("traffic_exporter"))
		os.Exit(0)
	}
	exporter := NewExporter(*scrapeURI, *insecure)
	prometheus.MustRegister(exporter)
	prometheus.MustRegister(version.NewCollector("traffic_exporter"))

	log.Infoln("Starting traffic_exporter", version.Info())
	log.Infoln("Build context", version.BuildContext())
	log.Infof("Starting Server: %s", *listeningAddress)

	http.Handle(*metricsEndpoint, prometheus.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
				<head><title>Apache Exporter</title></head>
				<body>
				<h1>Apache Exporter</h1>
				<p><a href='` + *metricsEndpoint + `'>Metrics</a></p>
				</body>
				</html>`))
	})
	log.Fatal(http.ListenAndServe(*listeningAddress, nil))
}
