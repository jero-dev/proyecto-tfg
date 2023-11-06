package memory

import (
	"testing"

	aggregates "vidya-sale/api/aggregate"
	"vidya-sale/api/domain/product"

	"github.com/google/uuid"
)

var (
	// testProduct is the name of the product used for testing
	testProduct = "Test Product"
	// testPlatform is the name of the platform used for testing
	testPlatform = "Test Platform"
	// testLink is the link used for testing
	testLink = "https://www.fake-shop.com"
	// testNoMatchingLink is the link used for testing when there is no matching link
	testNoMatchingLink = "https://www.fake-shop2.com"
)

func TestMemoryProductRepository_GetAll(t *testing.T) {
	repo := New()

	type testCase struct {
		name          string
		products      []aggregates.Product
		expectedError error
	}

	existingProduct, _ := aggregates.NewProduct(testProduct, testPlatform)

	testCases := []testCase{
		{
			name:          "No products in repository",
			products:      []aggregates.Product{},
			expectedError: product.ErrorThereAreNoProducts,
		},
		{
			name:          "Products in repository",
			products:      []aggregates.Product{existingProduct},
			expectedError: nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			for _, product := range testCase.products {
				repo.Add(product)
			}

			_, testError := repo.GetAll()
			if testError != testCase.expectedError {
				t.Errorf("Expected error to be %v but got %v", testCase.expectedError, testError)
			}
		})
	}
}

func TestMemorySaleRepository_GetByID(t *testing.T) {
	repo := New()

	type testCase struct {
		name          string
		products      []aggregates.Product
		id            uuid.UUID
		expectedError error
	}

	existingProduct, _ := aggregates.NewProduct(testProduct, testPlatform)

	testCases := []testCase{
		{
			name:          "No sales in repository",
			products:      []aggregates.Product{},
			id:            existingProduct.GetID(),
			expectedError: product.ErrorProductNotFound,
		},
		{
			name:          "Sales in repository but no matching link",
			products:      []aggregates.Product{existingProduct},
			id:            uuid.New(),
			expectedError: product.ErrorProductNotFound,
		},
		{
			name:          "Sales in repository and matching link",
			products:      []aggregates.Product{existingProduct},
			id:            existingProduct.GetID(),
			expectedError: nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			for _, product := range testCase.products {
				repo.Add(product)
			}

			_, testError := repo.GetByID(testCase.id)
			if testError != testCase.expectedError {
				t.Errorf("Expected error to be %v but got %v", testCase.expectedError, testError)
			}
		})
	}
}

func TestMemoryProductRepository_Add(t *testing.T) {
	repo := New()

	type testCase struct {
		name          string
		products      []aggregates.Product
		product       aggregates.Product
		expectedError error
	}

	existingProduct, _ := aggregates.NewProduct(testProduct, testPlatform)

	testCases := []testCase{
		{
			name:          "Product does not exist in repository",
			products:      []aggregates.Product{},
			product:       existingProduct,
			expectedError: nil,
		},
		{
			name:          "Sale exists in repository",
			products:      []aggregates.Product{existingProduct},
			product:       existingProduct,
			expectedError: product.ErrorFailedToAddProduct,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			for _, product := range testCase.products {
				repo.Add(product)
			}

			testError := repo.Add(testCase.product)
			if testError != testCase.expectedError {
				t.Errorf("Expected error to be %v but got %v", testCase.expectedError, testError)
			}
		})
	}
}

func TestMemoryProductRepository_Update(t *testing.T) {
	repo := New()

	type testCase struct {
		name          string
		products      []aggregates.Product
		product       aggregates.Product
		expectedError error
	}

	existingProduct, _ := aggregates.NewProduct(testProduct, testPlatform)

	testCases := []testCase{
		{
			name:          "Product does not exist in repository",
			products:      []aggregates.Product{},
			product:       existingProduct,
			expectedError: product.ErrorUpdateProduct,
		},
		{
			name:          "Product exists in repository",
			products:      []aggregates.Product{existingProduct},
			product:       existingProduct,
			expectedError: nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			for _, Product := range testCase.products {
				repo.Add(Product)
			}

			testError := repo.Update(testCase.product)
			if testError != testCase.expectedError {
				t.Errorf("Expected error to be %v but got %v", testCase.expectedError, testError)
			}
		})
	}
}
