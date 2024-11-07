package counters

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var Uptime = promauto.NewGauge(
	prometheus.GaugeOpts{
		Name: "mhz19_uptime",
		Help: "Mhz19 Application Uptime In Seconds",
	},
)

var ApiRequests = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Name: "mhz19_api_requests",
		Help: "Number Of Total Requests To Api Split By path And method",
	},
	[]string{"path", "method"},
)

var Transactions = promauto.NewHistogram(
	prometheus.HistogramOpts{
		Name: "mhz19_transactions",
		Help: "Mhz19 Count Of Transaction On Sqlite3 Db And Execution Time In Seconds",
	},
)

var BuildRules = promauto.NewHistogram(
	prometheus.HistogramOpts{
		Name: "mhz19_build_rules",
		Help: "Mhz19 time taken to transform rules from db representation into the types.Rule",
	},
)

var Queries = promauto.NewHistogram(
	prometheus.HistogramOpts{
		Name: "mhz19_queries",
		Help: "Mhz19 Count Of Queries On Sqlite3 Db And Execution Time In Seconds",
	},
)

var MessagesByChannel = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Name: "mhz19_messages_by_channel",
		Help: "Mhz19 Messages Handled Count Split By Channel",
	},
	[]string{"channel"},
)

var MessagesHandled = promauto.NewHistogram(
	prometheus.HistogramOpts{
		Name: "mhz19_messages_handled",
		Help: "Mhz19 Messages Handled Count And Execution Time In Seconds",
	},
)

var Errors = promauto.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "mhz19_errors",
		Help: "Mhz19 Errors Split By Module",
	},
	[]string{"module"},
)

var Rules = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Name: "mhz19_rules",
		Help: "Mhz19 Rules Executions Split By Rule Name",
	},
	[]string{"rule_name"},
)

// Buckets: prometheus.LinearBuckets(
// 	100,      // 100us - fastest expected
// 	500*1000, // 500ms - bucket size
// 	10,       // 10x500ms=5s - slowest expected, total buckets
// ),
