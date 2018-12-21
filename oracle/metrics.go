package oracle

import (
	"fmt"
	"time"

	"github.com/go-kit/kit/metrics"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	methodCallCount                metrics.Counter
	methodDuration                 metrics.Histogram
	fetchedPlasmachainEventCount   metrics.Counter
	submittedPlasmachainEventCount metrics.Counter
}

func NewMetrics(subsystem string) *Metrics {
	const namespace = "loomchain"

	return &Metrics{
		methodCallCount: kitprometheus.NewCounterFrom(
			stdprometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "method_call_count",
				Help:      "Number of times a method has been invoked.",
			}, []string{"method", "error"}),
		methodDuration: kitprometheus.NewSummaryFrom(
			stdprometheus.SummaryOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "method_duration",
				Help:      "How long a method took to execute (in seconds).",
			}, []string{"method", "error"}),
		fetchedPlasmachainEventCount: kitprometheus.NewCounterFrom(
			stdprometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "fetched_plasma_event_count",
				Help:      "Number of Plasmachain events fetched from the Plasmachain.",
			}, []string{"kind"}),
		submittedPlasmachainEventCount: kitprometheus.NewCounterFrom(
			stdprometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "submitted_plasma_event_count",
				Help:      "Number of Plasmachain events successfully submitted to the Gamechain.",
			}, nil),
	}
}

func (m *Metrics) MethodCalled(begin time.Time, method string, err error) {
	lvs := []string{"method", method, "error", fmt.Sprint(err != nil)}
	m.methodDuration.With(lvs...).Observe(time.Since(begin).Seconds())
	m.methodCallCount.With(lvs...).Add(1)
}

func (m *Metrics) FetchedMPlasmachainEvents(numEvents int, kind string) {
	m.fetchedPlasmachainEventCount.With("kind", kind).Add(float64(numEvents))
}

func (m *Metrics) SubmittedPlasmachainEvents(numEvents int) {
	m.submittedPlasmachainEventCount.Add(float64(numEvents))
}
