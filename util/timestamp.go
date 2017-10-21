package util

import (
	strftime "github.com/jehiah/go-strftime"
	"strings"
	"time"
)

// Timestamp returns the timestamp in Letto's standard format:
// `%Y%m%dT%H%M%S%L%Z`.
func Timestamp(t time.Time) string {
	timestamp := strftime.Format("%Y%m%dT%H%M%S%L%Z", t)
	return strings.Replace(timestamp, ".", "", 1)
}
