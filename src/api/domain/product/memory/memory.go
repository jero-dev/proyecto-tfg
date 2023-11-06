// Package memory is a in-memory implementation of the product repository
package memory

import (
	"sync"
	aggregates "vidya-sale/api/aggregate"
	product "vidya-sale/api/domain/product"

	"github.com/google/uuid"
)

// MemoryRepository fulfills the ProductRepository interface
type MemoryRepository struct {
	products map[uuid.UUID]aggregates.Product
	sync.Mutex
}

// New is a factory function to generate a new repository of products
func New() *MemoryRepository {
	return &MemoryRepository{
		products: make(map[uuid.UUID]aggregates.Product),
	}
}

// GetAll finds all the products that are in the repository
func (mr *MemoryRepository) GetAll() ([]aggregates.Product, error) {
	if length := len(mr.products); length == 0 {
		return []aggregates.Product{}, product.ErrorThereAreNoProducts
	}

	products := []aggregates.Product{}

	for _, product := range mr.products {
		products = append(products, product)
	}

	return products, nil
}

// Get finds a product by its ID
func (mr *MemoryRepository) GetByID(id uuid.UUID) (aggregates.Product, error) {
	for _, product := range mr.products {
		if product.GetID() == id {
			return product, nil
		}
	}

	return aggregates.Product{}, product.ErrorProductNotFound
}

// Add will add a new product to the repository
func (mr *MemoryRepository) Add(productToAdd aggregates.Product) error {
	if mr.products == nil {
		mr.Lock()
		mr.products = make(map[uuid.UUID]aggregates.Product)
		mr.Unlock()
	}
	// Make sure Product isn't already in the repository
	if _, ok := mr.products[productToAdd.GetID()]; ok {
		return product.ErrorFailedToAddProduct
	}
	mr.Lock()
	mr.products[productToAdd.GetID()] = productToAdd
	mr.Unlock()
	return nil
}

// Update will replace an existing product information with the new product information
func (mr *MemoryRepository) Update(updatedProduct aggregates.Product) error {
	// Make sure Product is in the repository
	if _, ok := mr.products[updatedProduct.GetID()]; !ok {
		return product.ErrorUpdateProduct
	}
	mr.Lock()
	mr.products[updatedProduct.GetID()] = updatedProduct
	mr.Unlock()
	return nil
}
