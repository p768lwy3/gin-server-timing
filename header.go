package servertiming

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/golang/gddo/httputil/header"
)

const HeaderKey = "Server-Timing"

type Header struct {
	Metrics []*Metric
	sync.Mutex
}

func ParseHeader(input string) (*Header, error) {
	rawMetrics := header.ParseList(headerParams(input))

	metrics := make([]*Metric, 0, len(rawMetrics))
	for _, raw := range rawMetrics {
		var m Metric
		m.Name, m.Extra = header.ParseValueAndParams(headerParams(raw))

		if v, ok := m.Extra[paramNameDesc]; ok {
			m.Desc = v
			delete(m.Extra, paramNameDesc)
		}

		if v, ok := m.Extra[paramNameDur]; ok {
			m.Duration, _ = time.ParseDuration(v + "ms")
			delete(m.Extra, paramNameDur)
		}

		metrics = append(metrics, &m)
	}

	return &Header{Metrics: metrics}, nil
}

func (h *Header) NewMetric(name string) *Metric {
	return h.Add(&Metric{Name: name})
}

func (h *Header) Add(m *Metric) *Metric {
	if h == nil {
		return m
	}

	h.Lock()
	defer h.Unlock()
	h.Metrics = append(h.Metrics, m)
	return m
}

func (h *Header) String() string {
	parts := make([]string, 0, len(h.Metrics))
	for _, m := range h.Metrics {
		parts = append(parts, m.String())
	}

	return strings.Join(parts, ",")
}

const (
	paramNameDesc = "desc"
	paramNameDur  = "dur"
)

func headerParams(s string) (http.Header, string) {
	const key = "Key"
	return http.Header(map[string][]string{
		key: {s},
	}), key
}

var reNumber = regexp.MustCompile(`^\d+\.?\d*$`)

func headerEncodeParam(key, value string) string {
	if reNumber.MatchString(value) {
		return fmt.Sprintf(`%s=%s`, key, value)
	}

	return fmt.Sprintf(`%s=%q`, key, value)
}
