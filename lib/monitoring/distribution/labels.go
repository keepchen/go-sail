package distribution

import (
	"crypto/sha256"
	"encoding/hex"
	"sort"
	"strings"
)

// Labels identifies a metric series. Callers should not modify the map after passing it
// to a Registry observation method.
type Labels map[string]string

// Clone returns an independent copy of the label set.
func (l Labels) Clone() Labels {
	result := make(Labels, len(l))
	for key, value := range l {
		result[key] = value
	}
	return result
}

// CanonicalString returns a stable label string sorted by key.
func (l Labels) CanonicalString() string {
	if len(l) == 0 {
		return ""
	}

	keys := make([]string, 0, len(l))
	for key := range l {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var builder strings.Builder
	for _, key := range keys {
		builder.WriteString(key)
		builder.WriteByte('=')
		builder.WriteString(l[key])
		builder.WriteByte('\n')
	}
	return builder.String()
}

// BuildMetricKey returns a stable identifier for a metric series.
func BuildMetricKey(namespace string, service string, metric string, labels Labels) string {
	raw := strings.Join([]string{
		namespace,
		service,
		metric,
		labels.CanonicalString(),
	}, "\x00")

	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}
