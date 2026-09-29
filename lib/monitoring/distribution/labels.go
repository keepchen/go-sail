package distribution

import (
	"crypto/sha256"
	"encoding/hex"
	"sort"
	"strings"
)

// Labels 标识一条指标序列。调用方将其传给 Registry 的观测方法后，
// 不应再修改该映射。
type Labels map[string]string

// Clone 返回标签集合的独立副本。
func (l Labels) Clone() Labels {
	result := make(Labels, len(l))
	for key, value := range l {
		result[key] = value
	}
	return result
}

// CanonicalString 返回按键排序且结果稳定的标签字符串。
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

// BuildMetricKey 返回指标序列的稳定标识。
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
