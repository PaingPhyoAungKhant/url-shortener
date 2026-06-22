// Package domain
package domain

import "time"

// Entry stores the mapping between a shortened code and the original URL.
type Entry struct {
	Code      string
	URL       string
	Visits    int64
	CreatedAt time.Time
}
