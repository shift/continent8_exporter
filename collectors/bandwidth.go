package collectors

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
)

// Network data structure
type Network struct {
	Time string `json:"time"`
	In   string `json:"in"`
	Out  string `json:"out"`
}

// BandwidthRack data structure
type BandwidthRack struct {
	Time     string `json:"time"`
	In       string `json:"in"`
	Out      string `json:"out"`
	Networks map[string]Network
}

// BandwidthDatacentre data structure
type BandwidthDatacentre map[string]BandwidthRack

// Bandwidth data structure
type Bandwidth map[string]BandwidthDatacentre

// BandwidthCollector definition
type BandwidthCollector struct {
	counterDesc *prometheus.Desc
}

// Describe definition
func (c *BandwidthCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.counterDesc
}

// Collect definition
func (c *BandwidthCollector) Collect(ch chan<- prometheus.Metric) {
	defer func() {
		if r := recover(); r != nil {
			level.Error(logger).Log("msg", "recovered", "err", r)
		}
	}()

	client := http.Client{
		Timeout: time.Second * 2,
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("User-Agent", "continent8_exporter/v0.0.1")

	q := req.URL.Query()
	q.Add("action", "getBandwidth")
	q.Add("username", os.Getenv("C8_USERNAME"))
	q.Add("token", os.Getenv("C8_TOKEN"))
	req.URL.RawQuery = q.Encode()

	start := time.Now()
	res, getErr := client.Do(req)
	if getErr != nil {
		panic(getErr)
	}

	duration := float64(time.Since(start).Seconds())
	scrapeDuration.WithLabelValues("bandwidth").Observe(duration)

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		panic(readErr)
	}
	bw := Bandwidth{}

	jsonErr := json.Unmarshal(body, &bw)
	if jsonErr != nil {
		panic(jsonErr)
	}

	for dc, v := range bw {
		for k, v := range v {
			inf, err := strconv.ParseFloat(v.In, 64)
			if err != nil {
				panic(err)
			}
			outf, err := strconv.ParseFloat(v.Out, 64)
			if err != nil {
				panic(err)
			}

			ch <- prometheus.MustNewConstMetric(c.counterDesc, prometheus.CounterValue, inf, dc, k, v.Time, "in", "total")
			ch <- prometheus.MustNewConstMetric(c.counterDesc, prometheus.CounterValue, outf, dc, k, v.Time, "out", "total")
			for network, v := range v.Networks {
				ch <- prometheus.MustNewConstMetric(c.counterDesc, prometheus.CounterValue, inf, dc, k, v.Time, "in", network)
				ch <- prometheus.MustNewConstMetric(c.counterDesc, prometheus.CounterValue, outf, dc, k, v.Time, "out", network)
			}
		}
	}
}

// NewBandwidthCollector definition
func NewBandwidthCollector() *BandwidthCollector {
	return &BandwidthCollector{
		counterDesc: prometheus.NewDesc("bandwidth", "C8 Bandwidth",
			[]string{"datacentre", "rack", "time", "type", "network"},
			nil),
	}
}
