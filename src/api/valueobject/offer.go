// Package valueobjects holds all the value objects that are shared across all subdomains
package valueobjects

import "errors"

var (
	// ErrorMissingProperties is returned when one of the obligatory parameters is missing.
	ErrorMissingProperties = errors.New("there is a missing property. You need to provide a link and a price")
)

// Offer represents an offer in the system
type Offer struct {
	link  string
	price float64
}

func NewOffer(link string, price float64) (Offer, error) {
	if link == "" || price == 0 {
		return Offer{}, ErrorMissingProperties
	}

	return Offer{
		link:  link,
		price: price,
	}, nil
}

func (offer Offer) GetLink() string {
	return offer.link
}

func (offer Offer) GetPrice() float64 {
	return offer.price
}
