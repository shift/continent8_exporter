package main

import (
	"net/http"
	"os"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/promlog"
	"github.com/prometheus/common/promlog/flag"
	"github.com/prometheus/common/version"
	"github.com/shift/continent8_exporter/collectors"
	"gopkg.in/alecthomas/kingpin.v2"
)

// TODO: use https://support.continent8.com/data/locations.json to prettify datacenter names
// TODO: use prettify rack names

var (
	httpBind       = kingpin.Flag("bind", "The address to listen on for HTTP requests.").Default(":9364").String()
	scrapeDuration = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name: "scrape_collection_duration_seconds",
			Help: "Duration of scraping",
		},
		[]string{"method"},
	)
	logger = log.NewNopLogger()
)

func init() {
	prometheus.MustRegister(version.NewCollector("continent8_exporter"))
	prometheus.MustRegister(scrapeDuration)
}
func main() {
	allowedLevel := promlog.AllowedLevel{}
	flag.AddFlags(kingpin.CommandLine, &allowedLevel)
	kingpin.Version(version.Print("continent8_exporter"))
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()
	logger = promlog.New(allowedLevel)
	collectors.Logger(logger)
	level.Info(logger).Log("msg", "Starting continent8_exporter", "version", version.Info())
	level.Info(logger).Log("msg", "Build context", version.BuildContext())
	prometheus.MustRegister(collectors.NewBandwidthCollector())
	prometheus.MustRegister(collectors.NewEnvironmentCollector())
	http.Handle("/metrics", promhttp.Handler())
	level.Info(logger).Log("msg", "Listening", "port", *httpBind)
	if err := http.ListenAndServe(*httpBind, nil); err != nil {
		level.Error(logger).Log("msg", "Error starting HTTP server", "err", err)
		os.Exit(1)
	}
}
