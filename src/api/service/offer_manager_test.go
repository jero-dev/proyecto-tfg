package services

import (
	"testing"
	aggregates "vidya-sale/api/aggregate"
	valueobjects "vidya-sale/api/valueobject"
)

func TestOfferManagerService_StoreOffer(t *testing.T) {
	existingProduct, _ := aggregates.NewProduct("Test Game Name", "Test Game Platform")
	existingProductWithOffer, _ := aggregates.NewProduct("Test Game Name", "Test Game Platform")
	existingOffer, _ := valueobjects.NewOffer("Test Offer Link", 5.0)
	existingProductWithOffer.AddOffer(existingOffer)
	expectedProduct, _ := aggregates.NewProduct("Test Game Name", "Test Game Platform")
	expectedOffer, _ := valueobjects.NewOffer("Test Offer Link", 10.0)
	expectedProduct.AddOffer(expectedOffer)
	expectedProductWithTwoOffers, _ := aggregates.NewProduct("Test Game Name", "Test Game Platform")
	expectedOffer2, _ := valueobjects.NewOffer("Test Offer Link 2", 20.0)
	expectedProductWithTwoOffers.AddOffer(expectedOffer)
	expectedProductWithTwoOffers.AddOffer(expectedOffer2)

	type testCase struct {
		name             string
		gameName         string
		gamePlatform     string
		offerLink        string
		offerPrice       float64
		existingProducts []aggregates.Product
		expectedData     []aggregates.Product
		expectedError    error
	}

	testCases := []testCase{
		{
			name:             "New Product and Offer",
			gameName:         "Test Game Name",
			gamePlatform:     "Test Game Platform",
			offerLink:        "Test Offer Link",
			offerPrice:       10.0,
			existingProducts: []aggregates.Product{},
			expectedData: []aggregates.Product{
				expectedProduct,
			},
			expectedError: nil,
		},
		{
			name:         "Existing Product and New Offer",
			gameName:     "Test Game Name",
			gamePlatform: "Test Game Platform",
			offerLink:    "Test Offer Link",
			offerPrice:   10.0,
			existingProducts: []aggregates.Product{
				existingProduct,
			},
			expectedData: []aggregates.Product{
				expectedProduct,
			},
			expectedError: nil,
		},
		{
			name:         "Existing Product and New Offer on same link",
			gameName:     "Test Game Name",
			gamePlatform: "Test Game Platform",
			offerLink:    "Test Offer Link",
			offerPrice:   10.0,
			existingProducts: []aggregates.Product{
				existingProductWithOffer,
			},
			expectedData: []aggregates.Product{
				expectedProduct,
			},
			expectedError: nil,
		},

		{
			name:         "Existing Product and New Offer on other link",
			gameName:     "Test Game Name",
			gamePlatform: "Test Game Platform",
			offerLink:    "Test Offer Link 2",
			offerPrice:   20.0,
			existingProducts: []aggregates.Product{
				existingProductWithOffer,
			},
			expectedData: []aggregates.Product{
				expectedProductWithTwoOffers,
			},
			expectedError: nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			managerService, _ := NewOfferManagerService(WithMemoryProductRepository())
			for _, existingProduct := range testCase.existingProducts {
				managerService.products.Add(existingProduct)
			}

			testError := managerService.StoreOffer(testCase.gameName, testCase.gamePlatform, testCase.offerLink, testCase.offerPrice)
			actualProducts, _ := managerService.products.GetAll()

			if testError != testCase.expectedError {
				t.Errorf("Expected error to be %v but got %v", testCase.expectedError, testError)
			}
			if len(actualProducts) != len(testCase.expectedData) {
				t.Errorf("Expected %v products but got %v", len(testCase.expectedData), len(actualProducts))
			}
			if actualProducts[0].GetName() != testCase.expectedData[0].GetName() {
				t.Errorf("Expected product name to be %v but got %v", testCase.expectedData[0].GetName(), actualProducts[0].GetName())
			}
			if actualProducts[0].GetPlatform() != testCase.expectedData[0].GetPlatform() {
				t.Errorf("Expected product platform to be %v but got %v", testCase.expectedData[0].GetPlatform(), actualProducts[0].GetPlatform())
			}
			if len(actualProducts[0].GetOffers()) != len(testCase.expectedData[0].GetOffers()) {
				t.Errorf("Expected %v offers but got %v", len(testCase.expectedData[0].GetOffers()), len(actualProducts[0].GetOffers()))
			}
			for index, actualOffer := range actualProducts[0].GetOffers() {
				if actualOffer.GetLink() != testCase.expectedData[0].GetOffers()[index].GetLink() {
					t.Errorf("Expected offer link to be %v but got %v", testCase.expectedData[0].GetOffers()[0].GetLink(), actualOffer.GetLink())
				}
				if actualOffer.GetPrice() != testCase.expectedData[0].GetOffers()[index].GetPrice() {
					t.Errorf("Expected offer price to be %v but got %v", testCase.expectedData[0].GetOffers()[0].GetPrice(), actualOffer.GetPrice())
				}
			}
		})
	}
}

func TestOfferManagerService_GetGameOffers(t *testing.T) {
	existingOffer, _ := valueobjects.NewOffer("Test Offer Link", 5.0)
	existingOffer2, _ := valueobjects.NewOffer("Test Offer Link 2", 20.0)
	existingProductWithOffer, _ := aggregates.NewProduct("Test Game Name", "Test Game Platform")
	existingProductWithOffer.AddOffer(existingOffer)
	existingProductInAnotherPlatform, _ := aggregates.NewProduct("Test Game Name", "Another Test Game Platform")
	existingProductInAnotherPlatform.AddOffer(existingOffer)
	existingProductInAnotherPlatform.AddOffer(existingOffer2)
	existingProductWithTwoOffers, _ := aggregates.NewProduct("Another Test Game Name", "Test Game Platform")
	existingProductWithTwoOffers.AddOffer(existingOffer)
	existingProductWithTwoOffers.AddOffer(existingOffer2)
	expectedOffers := make(map[string][]valueobjects.Offer)
	expectedOffers["Test Game Platform"] = []valueobjects.Offer{existingOffer}
	expectedOffers["Another Test Game Platform"] = []valueobjects.Offer{existingOffer, existingOffer2}

	type testCase struct {
		name             string
		gameName         string
		existingProducts []aggregates.Product
		expectedData     map[string][]valueobjects.Offer
		expectedError    error
	}

	testCases := []testCase{
		{
			name:             "No products in the repository",
			gameName:         "Test Game Name",
			existingProducts: []aggregates.Product{},
			expectedData:     make(map[string][]valueobjects.Offer),
			expectedError:    ErrorThereAreNoOffers,
		},
		{
			name:     "No offers for the given game name",
			gameName: "Random Game Name",
			existingProducts: []aggregates.Product{
				existingProductWithOffer, existingProductInAnotherPlatform, existingProductWithTwoOffers,
			},
			expectedData:  make(map[string][]valueobjects.Offer),
			expectedError: ErrorThereAreNoOffers,
		},
		{
			name:     "Existing offers for the given game name",
			gameName: "Test Game Name",
			existingProducts: []aggregates.Product{
				existingProductWithOffer, existingProductInAnotherPlatform, existingProductWithTwoOffers,
			},
			expectedData:  expectedOffers,
			expectedError: nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			managerService, _ := NewOfferManagerService(WithMemoryProductRepository())
			for _, existingProduct := range testCase.existingProducts {
				managerService.products.Add(existingProduct)
			}

			gameOffers, testError := managerService.GetGameOffers(testCase.gameName)

			if testError != testCase.expectedError {
				t.Errorf("Expected error to be %v but got %v", testCase.expectedError, testError)
			}
			if len(gameOffers) != len(testCase.expectedData) {
				t.Errorf("Expected %v products but got %v", len(testCase.expectedData), len(gameOffers))
			}
			for platform, offers := range gameOffers {
				if testCase.expectedData[platform] == nil {
					t.Errorf("Expected platform %v to have offers", platform)
				}
				if len(offers) != len(testCase.expectedData[platform]) {
					t.Errorf("Expected %v offers but got %v", len(testCase.expectedData[platform]), len(offers))
				}
				for index, offer := range offers {
					if offer.GetLink() != testCase.expectedData[platform][index].GetLink() {
						t.Errorf("Expected offer link to be %v but got %v", testCase.expectedData[platform][index].GetLink(), offer.GetLink())
					}
					if offer.GetPrice() != testCase.expectedData[platform][index].GetPrice() {
						t.Errorf("Expected offer price to be %v but got %v", testCase.expectedData[platform][index].GetPrice(), offer.GetPrice())
					}
				}
			}
		})
	}
}
