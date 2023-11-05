// Package services holds all the services that connects repositories into a business flow
package services

import aggregates "vidya-sale/api/aggregate"

// OfferManager is an interface that defines the rules around what an offer manager
// service has to be able to perform
type OfferManager interface {
	StoreOffer(message string) error
	GetProductOffers(productName string) ([]aggregates.Product, error)
}
