package prometheusMetrics

import "github.com/prometheus/client_golang/prometheus"

var BookAddCounter = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "book_add_count",
		Help: "Number of books successfully added to database",
	},
)

var BookGetounter = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "book_get_count",
		Help: "Number of get requests successfully completed",
	},
)
