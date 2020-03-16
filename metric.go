package servertiming

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Metric struct {
	Name      string
	Duration  time.Duration
	Desc      string
	Extra     map[string]string
	startTime time.Time
}

func (m *Metric) WithDesc(desc string) *Metric {
	m.Desc = desc
	return m
}

func (m *Metric) Start() *Metric {
	m.startTime = time.Now()
	return m
}

func (m *Metric) Stop() *Metric {
	if !m.startTime.IsZero() {
		m.Duration = time.Since(m.startTime)
	}

	return m
}

func (m *Metric) String() string {

	parts := make([]string, 1, len(m.Extra)+3)
	parts[0] = m.Name

	if _, ok := m.Extra[paramNameDesc]; !ok && m.Desc != "" {
		parts = append(parts, headerEncodeParam(paramNameDesc, m.Desc))
	}

	if _, ok := m.Extra[paramNameDur]; !ok && m.Duration > 0 {
		parts = append(parts, headerEncodeParam(
			paramNameDur,
			strconv.FormatFloat(float64(m.Duration)/float64(time.Millisecond), 'f', -1, 64),
		))
	}

	for k, v := range m.Extra {
		parts = append(parts, headerEncodeParam(k, v))
	}

	return strings.Join(parts, ";")
}

func (m *Metric) GoString() string {
	if m == nil {
		return "nil"
	}

	return fmt.Sprintf("*%#v", *m)
}
