package prometheusMetrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	prometheus.MustRegister(BookGetCounter, BookAddCounter, BookDeleteCounter, BookGetAllCounter, BookUpdateCounter)
}

var BookAddCounter = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "book_add_count",
		Help: "Number of books successfully added to database",
	},
)

var BookGetCounter = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "book_get_count",
		Help: "Number of get requests successfully completed",
	},
)

var BookGetAllCounter = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "book_get_all_count",
		Help: "Number of \"get all books\" requests successfully completed",
	},
)

var BookDeleteCounter = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "book_delete_count",
		Help: "Number of delete successfully completed",
	},
)

var BookUpdateCounter = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "book_update_count",
		Help: "Number of update requests successfully completed",
	},
)
