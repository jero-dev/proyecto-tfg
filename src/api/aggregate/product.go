// Package aggregates holds all the aggregates that are shared across all subdomains
package aggregates

import (
	"errors"
	entities "vidya-sale/api/entity"
	valueobjects "vidya-sale/api/valueobject"

	"github.com/google/uuid"
)

var (
	// ErrorMissingProperties is returned when one of the obligatory parameters is missing.
	ErrorMissingProperties = errors.New("there is a missing property. You need to provide a name and a platform")
)

// Product represents a product in the system
type Product struct {
	videoGame *entities.VideoGame
	offers    []valueobjects.Offer
}

func NewProduct(name, platform string) (Product, error) {

	if name == "" || platform == "" {
		return Product{}, ErrorMissingProperties
	}

	videoGame := &entities.VideoGame{
		ID:       uuid.New(),
		Name:     name,
		Platform: platform,
	}

	return Product{
		videoGame: videoGame,
		offers:    make([]valueobjects.Offer, 0),
	}, nil
}

func (product *Product) GetID() uuid.UUID {
	return product.videoGame.ID
}

func (product *Product) GetName() string {
	return product.videoGame.Name
}

func (product *Product) GetPlatform() string {
	return product.videoGame.Platform
}

func (product *Product) GetOffers() []valueobjects.Offer {
	return product.offers
}

func (product *Product) AddOffer(offer valueobjects.Offer) {
	for index, currentOffer := range product.offers {
		if currentOffer.GetLink() == offer.GetLink() {
			product.offers = append(product.offers[:index], product.offers[index+1:]...)
		}
	}
	product.offers = append(product.offers, offer)
}
