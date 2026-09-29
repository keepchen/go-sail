package redisbackend

import (
	"fmt"
	"time"
)

const prefix = "go-sail:distribution"

func recordKey(recordID string) string {
	return fmt.Sprintf("%s:record:%s", prefix, recordID)
}

func windowKey(namespace string, service string, windowStart time.Time) string {
	return fmt.Sprintf(
		"%s:window:%s:%s:%d",
		prefix,
		namespace,
		service,
		windowStart.Unix(),
	)
}

func resultKey(resultID string) string {
	return fmt.Sprintf("%s:result:%s", prefix, resultID)
}

func aggregatedKey(recordID string) string {
	return fmt.Sprintf("%s:aggregated:%s", prefix, recordID)
}

func cursorKey(namespace string, service string) string {
	return fmt.Sprintf("%s:cursor:%s:%s", prefix, namespace, service)
}
