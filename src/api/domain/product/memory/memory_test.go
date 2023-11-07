package memory_test

import (
	"testing"

	aggregates "vidya-sale/api/aggregate"
	"vidya-sale/api/domain/product"
	"vidya-sale/api/domain/product/memory"

	"github.com/google/uuid"
)

var (
	// testVideoGameName is the name of the product used for testing
	testVideoGameName = "Test Video Game Name"
	// testVideoGamePlatform is the name of the platform used for testing
	testVideoGamePlatform = "Test Video Game Platform"
)

func TestMemoryProductRepository_GetAll(t *testing.T) {
	repo := memory.New()

	type testCase struct {
		name          string
		products      []aggregates.Product
		expectedError error
	}

	existingProduct, _ := aggregates.NewProduct(testVideoGameName, testVideoGamePlatform)

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
	repo := memory.New()

	type testCase struct {
		name          string
		products      []aggregates.Product
		id            uuid.UUID
		expectedError error
	}

	existingProduct, _ := aggregates.NewProduct(testVideoGameName, testVideoGamePlatform)

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

func TestMemoryProductRepository_GetByNameAndPlatform(t *testing.T) {
	repo := memory.New()

	type testCase struct {
		name              string
		products          []aggregates.Product
		videoGameName     string
		videoGamePlatform string
		expectedError     error
	}

	existingProduct, _ := aggregates.NewProduct(testVideoGameName, testVideoGamePlatform)

	testCases := []testCase{
		{
			name:              "No products in repository",
			products:          []aggregates.Product{},
			videoGameName:     testVideoGameName,
			videoGamePlatform: testVideoGamePlatform,
			expectedError:     product.ErrorProductNotFound,
		},
		{
			name:              "Products in repository but no matching name and platform",
			products:          []aggregates.Product{existingProduct},
			videoGameName:     "Not Test Product",
			videoGamePlatform: "Not Test Platform",
			expectedError:     product.ErrorProductNotFound,
		},
		{
			name:              "Products in repository but no matching name",
			products:          []aggregates.Product{existingProduct},
			videoGameName:     "Not Test Product",
			videoGamePlatform: testVideoGamePlatform,
			expectedError:     product.ErrorProductNotFound,
		},
		{
			name:              "Products in repository but no matching platform",
			products:          []aggregates.Product{existingProduct},
			videoGameName:     testVideoGameName,
			videoGamePlatform: "Not Test Platform",
			expectedError:     product.ErrorProductNotFound,
		},
		{
			name:              "Products in repository and matching name and platform",
			products:          []aggregates.Product{existingProduct},
			videoGameName:     testVideoGameName,
			videoGamePlatform: testVideoGamePlatform,
			expectedError:     nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			for _, product := range testCase.products {
				repo.Add(product)
			}

			_, testError := repo.GetByNameAndPlatform(testCase.videoGameName, testCase.videoGamePlatform)
			if testError != testCase.expectedError {
				t.Errorf("Expected error to be %v but got %v", testCase.expectedError, testError)
			}
		})
	}
}

func TestMemoryProductRepository_GetByName(t *testing.T) {
	repo := memory.New()

	type testCase struct {
		name             string
		products         []aggregates.Product
		videoGameName    string
		expectedProducts []aggregates.Product
		expectedError    error
	}

	existingProduct, _ := aggregates.NewProduct(testVideoGameName, testVideoGamePlatform)
	existingProductInDifferentPlatform, _ := aggregates.NewProduct(testVideoGameName, "Not Test Platform")
	existingProductWithDifferentName, _ := aggregates.NewProduct("Not Test Product", testVideoGamePlatform)

	testCases := []testCase{
		{
			name:             "No products in repository",
			products:         []aggregates.Product{},
			videoGameName:    testVideoGameName,
			expectedProducts: []aggregates.Product{},
			expectedError:    product.ErrorProductNotFound,
		},
		{
			name:             "Products in repository but no matching name",
			products:         []aggregates.Product{existingProduct, existingProductInDifferentPlatform, existingProductWithDifferentName},
			videoGameName:    "Test Product Name",
			expectedProducts: []aggregates.Product{},
			expectedError:    product.ErrorProductNotFound,
		},
		{
			name:             "Products in repository and matching name",
			products:         []aggregates.Product{existingProduct, existingProductInDifferentPlatform, existingProductWithDifferentName},
			videoGameName:    testVideoGameName,
			expectedProducts: []aggregates.Product{existingProduct, existingProductInDifferentPlatform},
			expectedError:    nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			for _, product := range testCase.products {
				repo.Add(product)
			}

			testProducts, testError := repo.GetByName(testCase.videoGameName)
			if testError != testCase.expectedError {
				t.Errorf("Expected error to be %v but got %v", testCase.expectedError, testError)
			}
			if len(testProducts) != len(testCase.expectedProducts) {
				t.Errorf("Expected products to be %v but got %v", testCase.expectedProducts, testProducts)
			}
			for index, testProduct := range testProducts {
				if testProduct.GetName() != testCase.expectedProducts[index].GetName() {
					t.Errorf("Expected product name to be %v but got %v",
						testCase.expectedProducts[index].GetName(), testProduct.GetName())
				}
				if testProduct.GetPlatform() != testCase.expectedProducts[index].GetPlatform() {
					t.Errorf("Expected product platform to be %v but got %v",
						testCase.expectedProducts[index].GetPlatform(), testProduct.GetPlatform())
				}
			}
		})
	}
}

func TestMemoryProductRepository_Add(t *testing.T) {
	repo := memory.New()

	type testCase struct {
		name          string
		products      []aggregates.Product
		product       aggregates.Product
		expectedError error
	}

	existingProduct, _ := aggregates.NewProduct(testVideoGameName, testVideoGamePlatform)

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
	repo := memory.New()

	type testCase struct {
		name          string
		products      []aggregates.Product
		product       aggregates.Product
		expectedError error
	}

	existingProduct, _ := aggregates.NewProduct(testVideoGameName, testVideoGamePlatform)

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
