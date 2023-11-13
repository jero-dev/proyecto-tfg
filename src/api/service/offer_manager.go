// Package services holds all the services that connects repositories into a business flow
package services

import (
	"errors"
	aggregates "vidya-sale/api/aggregate"
	"vidya-sale/api/domain/product"
	"vidya-sale/api/domain/product/memory"
	"vidya-sale/api/domain/product/postgres"
	valueobjects "vidya-sale/api/valueobject"
)

var (
	// ErrorThereAreNoOffers is returned when there are no offers retrieve.
	ErrorThereAreNoOffers = errors.New("there are no offers for the given game")
)

type OfferManager interface {
	AddOffer(gameName, gamePlatform, offerLink string, offerPrice float64) error
	GetGameOffers(gameName string) ([]aggregates.Product, error)
}

// OfferManagerConfiguration is an alias for a function that will take in a pointer to
// a OfferManagerService and modify it
type OfferManagerConfiguration func(managerService *OfferManagerService) error

// OfferManagerService is an implementation of the OfferManagerService
type OfferManagerService struct {
	products product.ProductRepository
}

// NewOfferManagerService takes a variable amount of OfferManagerConfiguration
// functions and returns a new OfferManagerService.
// Each OfferManagerConfiguration will be called in the order they are passed in
func NewOfferManagerService(configurations ...OfferManagerConfiguration) (*OfferManagerService, error) {

	managerService := &OfferManagerService{}

	for _, configuration := range configurations {
		err := configuration(managerService)
		if err != nil {
			return nil, err
		}
	}

	return managerService, nil
}

// WithProductRepository applies a given product repository to the OfferManagerService
func WithProductRepository(productRepository product.ProductRepository) OfferManagerConfiguration {

	return func(processorService *OfferManagerService) error {
		processorService.products = productRepository
		return nil
	}
}

// WithPostgresProductRepository applies a Postgres product repository to the OfferManagerService
func WithPostgresProductRepository() OfferManagerConfiguration {
	productRepository, databaseError := postgres.New()
	if databaseError != nil {
		return func(processorService *OfferManagerService) error {
			return databaseError
		}
	}
	return WithProductRepository(productRepository)
}

// WithMemoryProductRepository applies a memory product repository to the OfferManagerService
func WithMemoryProductRepository() OfferManagerConfiguration {
	productRepository := memory.New()
	return WithProductRepository(productRepository)
}

// StoreOffer gets a message and stores the offer announced in the repository
func (managerService *OfferManagerService) StoreOffer(gameName, gamePlatform, offerLink string, offerPrice float64) error {
	foundProduct, findingError := managerService.products.GetByNameAndPlatform(gameName, gamePlatform)

	if findingError != nil {
		addingError := managerService.addProduct(gameName, gamePlatform, offerLink, offerPrice)

		return addingError
	} else {
		updatingError := managerService.updateProduct(foundProduct, offerLink, offerPrice)

		return updatingError
	}
}

// GetGameOffers gets a video game name and returns all the offers for that video game with all the respective platforms
func (managerService *OfferManagerService) GetGameOffers(gameName string) (map[string][]valueobjects.Offer, error) {
	foundProducts, findingError := managerService.products.GetByName(gameName)

	if findingError != nil {
		return nil, ErrorThereAreNoOffers
	}

	var platformOffers = make(map[string][]valueobjects.Offer)

	for _, product := range foundProducts {
		platformOffers[product.GetPlatform()] = product.GetOffers()
	}

	return platformOffers, nil
}

// addProduct gets a name, a platform, a link and a price and adds the resulting product to the repository
func (managerService *OfferManagerService) addProduct(name, platform, link string, price float64) error {
	newProduct, createProductError := aggregates.NewProduct(name, platform)

	if createProductError != nil {
		return createProductError
	}

	newOffer, createOfferError := valueobjects.NewOffer(link, price)

	if createOfferError != nil {
		return createOfferError
	}

	newProduct.AddOffer(newOffer)
	addProductError := managerService.products.Add(newProduct)

	if addProductError != nil {
		return addProductError
	}

	return nil
}

// updateProduct gets a product, a link and a price and updates it in the repository
func (managerService *OfferManagerService) updateProduct(product aggregates.Product, link string, price float64) error {
	newOffer, createOfferError := valueobjects.NewOffer(link, price)

	if createOfferError != nil {
		return createOfferError
	}

	product.AddOffer(newOffer)
	updateProductError := managerService.products.Update(product)

	if updateProductError != nil {
		return updateProductError
	}

	return nil
}
