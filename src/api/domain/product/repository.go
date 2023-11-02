package product

import (
	"errors"
	aggregates "vidya-sale/api/aggregate"

	"github.com/google/uuid"
)

var (
	// ErrorThereAreNoProducts is returned when there are no products in the repository.
	ErrorThereAreNoProducts = errors.New("there are no products in the repository")
	// ErrorProductNotFound is returned when a product is not found.
	ErrorProductNotFound = errors.New("the product was not found in the repository")
	// ErrorFailedToAddProduct is returned when the product could not be added to the repository.
	ErrorFailedToAddProduct = errors.New("failed to add the product to the repository")
	// ErrorUpdateProduct is returned when the product could not be updated in the repository.
	ErrorUpdateProduct = errors.New("failed to update the product in the repository")
)

// ProductRepository is a interface that defines the rules around what a product
// repository has to be able to perform
type ProductRepository interface {
	GetAll() ([]aggregates.Product, error)
	GetByID(uuid.UUID) (aggregates.Product, error)
	Add(aggregates.Product) error
	Update(aggregates.Product) error
}
