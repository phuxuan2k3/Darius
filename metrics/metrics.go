package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var LLMRequestCounter = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "llm_request_counter",
		Help: "The number of LLM requests made",
	},
	[]string{"feature"},
)

func init() {
	prometheus.MustRegister(LLMRequestCounter)
}
