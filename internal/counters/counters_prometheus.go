package counters

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var ApiRequests = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Name: "mhz19_api_requests_total",
		Help: "Number of api requests",
	},
	[]string{"path", "method"},
)

var MessagesHandled = promauto.NewCounter(
	prometheus.CounterOpts{
		Name: "mhz19_messages_handled",
	},
)

var Transactions = promauto.NewHistogram(
	prometheus.HistogramOpts{
		Name: "mhz19_transactions_ms",
	},
)
