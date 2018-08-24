package collectors

import (
	"github.com/go-kit/kit/log"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	scrapeDuration = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name: "scrape_collection_duration_seconds",
			Help: "Duration of scraping",
		},
		[]string{"method"},
	)

	url    = "https://support.continent8.com/api/"
	logger = log.NewNopLogger()
)

func Logger(plogger log.Logger) {
	logger = plogger
}
