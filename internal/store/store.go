// Package store
package store

import "github.com/PaingPhyoAungKhant/url-shortener/internal/domain"

type Store interface {
	// * Save the entry
	Save(originalURL string) (code string, err error)
	// * Get entry by code
	Get(code string) (entry domain.Entry, err error)
	// * Get code by url
	GetByURL(url string) (code string, err error)
	// * Incement visit count from code
	IncrementVisits(code string) error
	// * Delete code
	Delete(code string) (entry domain.Entry, err error)
}
