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
var LLMTokenCounter = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "llm_token_counter",
		Help: "The number of tokens used in LLM requests",
	},
	[]string{"feature"},
)

func init() {
	prometheus.MustRegister(LLMRequestCounter)
	prometheus.MustRegister(LLMTokenCounter)
}
