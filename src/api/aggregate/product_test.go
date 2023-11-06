package aggregates_test

import (
	"testing"

	aggregates "vidya-sale/api/aggregate"
	valueobjects "vidya-sale/api/valueobject"
)

func TestProduct_NewProduct(t *testing.T) {
	type TestCase struct {
		test          string
		product       string
		platform      string
		expectedError error
	}

	testCases := []TestCase{
		{
			test:          "Product is not provided",
			product:       "",
			platform:      "Test Platform",
			expectedError: aggregates.ErrorMissingProperties,
		},
		{
			test:          "Platform is not provided",
			product:       "Test Product",
			platform:      "",
			expectedError: aggregates.ErrorMissingProperties,
		},
		{
			test:          "Everything is provided",
			product:       "Test Product",
			platform:      "Test Platform",
			expectedError: nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.test, func(t *testing.T) {
			_, testError := aggregates.NewProduct(testCase.product, testCase.platform)
			if testError != testCase.expectedError {
				t.Errorf("Expected error to be %v but got %v", testCase.expectedError, testError)
			}
		})
	}
}

func TestProduct_AddOffer(t *testing.T) {
	type TestCase struct {
		test           string
		offer          valueobjects.Offer
		productOffers  []valueobjects.Offer
		expectedOffers []valueobjects.Offer
	}

	offer, _ := valueobjects.NewOffer("Test Link", 10.0)
	existingOffer, _ := valueobjects.NewOffer("Test Link", 5.0)

	testCases := []TestCase{
		{
			test:          "Offer does not exist",
			offer:         offer,
			productOffers: []valueobjects.Offer{},
			expectedOffers: []valueobjects.Offer{
				offer,
			},
		},
		{
			test:  "Offer already exists",
			offer: offer,
			productOffers: []valueobjects.Offer{
				existingOffer,
			},
			expectedOffers: []valueobjects.Offer{
				offer,
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.test, func(t *testing.T) {
			product, _ := aggregates.NewProduct("Test Product", "Test Platform")
			for _, offer := range testCase.productOffers {
				product.AddOffer(offer)
			}
			product.AddOffer(testCase.offer)
			offers := product.GetOffers()
			if offers[0].GetLink() != testCase.expectedOffers[0].GetLink() &&
				offers[0].GetPrice() != testCase.expectedOffers[0].GetPrice() {
				t.Errorf("Expected offer to be %v but got %v", testCase.expectedOffers, offers)
			}
		})
	}
}
