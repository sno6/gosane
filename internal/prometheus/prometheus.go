package prometheus

import "github.com/prometheus/client_golang/prometheus"

const (
	// Per generic request.
	CtxNameHandler = "handler"
	CtxNameMethod  = "method"
)

type Prometheus struct {
	HTTPRequestCounter   *prometheus.CounterVec
	HTTPResponseDuration *prometheus.CounterVec
}

func New() (*Prometheus, error) {
	p := &Prometheus{
		HTTPRequestCounter: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "http_request_total",
		}, []string{CtxNameHandler, CtxNameMethod}),
		HTTPResponseDuration: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "http_response_duration",
		}, []string{CtxNameHandler, CtxNameMethod}),
	}

	err := p.registerMetrics()
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (p *Prometheus) registerMetrics() error {
	for _, c := range []prometheus.Collector{
		p.HTTPRequestCounter,
		p.HTTPResponseDuration,
	} {
		err := prometheus.Register(c)
		if err != nil {
			return err
		}
	}

	return nil
}
