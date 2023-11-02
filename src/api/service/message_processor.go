// Package services holds all the services that connects repositories into a business flow
package services

import (
	"vidya-sale/api/domain/product"
)

// InfoProcessorService is an implementation of the service
type InfoProcessorService struct {
	products product.ProductRepository
}
