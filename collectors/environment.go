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

type Metric struct {
	Time  interface{} `json:"time"`  // N/A
	Value interface{} `json:"value"` // N/A
}
type EnvironmentMetrics map[string]Metric

// EnvironmentRack data structure
type EnvironmentRack map[string]Metric

// EnvironmentDatacentre data structure
type EnvironmentDatacentre map[string]EnvironmentRack

// Environment data structure
type Environment map[string]EnvironmentDatacentre

// EnvironmentCollector definition
type EnvironmentCollector struct {
	counterDesc *prometheus.Desc
}

// Describe definition
func (c *EnvironmentCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.counterDesc
}

// Collect definition
func (c *EnvironmentCollector) Collect(ch chan<- prometheus.Metric) {
	defer func() {
		if r := recover(); r != nil {
			// TODO: hide token from message
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
	q.Add("action", "getEnvironment")
	q.Add("username", os.Getenv("C8_USERNAME"))
	q.Add("token", os.Getenv("C8_TOKEN"))
	req.URL.RawQuery = q.Encode()

	start := time.Now()
	res, getErr := client.Do(req)
	if getErr != nil {
		panic(getErr)
	}

	duration := float64(time.Since(start).Seconds())
	scrapeDuration.WithLabelValues("environment").Observe(duration)

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		panic(readErr)
	}
	env := Environment{}

	jsonErr := json.Unmarshal(body, &env)
	if jsonErr != nil {
		panic(jsonErr)
	}

	for dc, racks := range env {
		for rack, rack_metrics := range racks {
			for metric_type, metric_values := range rack_metrics {
				value := 0.0
				switch metric_values.Value.(type) {
					case string:
						str := metric_values.Value.(string)
						if str != "N/A" {
							value, err = strconv.ParseFloat(str, 64)

							if err != nil {
								panic(err)
							}
						}
					case float64:
						value = metric_values.Value.(float64)
				}
				ch <- prometheus.MustNewConstMetric(c.counterDesc, prometheus.GaugeValue, value, dc, rack, metric_type)
			}
		}
	}
}

// NewEnvironmentCollector definition
func NewEnvironmentCollector() *EnvironmentCollector {
	return &EnvironmentCollector{
		counterDesc: prometheus.NewDesc("environment", "C8 Environment",
			[]string{"datacentre", "rack", "type"},
			nil),
	}
}
