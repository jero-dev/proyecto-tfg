package postgres

import (
	"log"
	"os"
	aggregates "vidya-sale/api/aggregate"
	"vidya-sale/api/domain/product"
	entities "vidya-sale/api/entity"
	valueobjects "vidya-sale/api/valueobject"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresRepository struct {
	database *gorm.DB
}

// New is a factory function to generate a new repository of products
func New() (*PostgresRepository, error) {
	var connectionError error

	database, connectionError := gorm.Open(postgres.Open(os.Getenv("DATABASE_DSN")), &gorm.Config{})

	if connectionError != nil {
		log.Print(connectionError.Error())
		return nil, connectionError
	}

	database.AutoMigrate(&entities.VideoGame{}, &valueobjects.Offer{}, &aggregates.Product{})

	return &PostgresRepository{
		database: database,
	}, nil
}

// GetAll finds all the products that are in the repository
func (postgresRepo *PostgresRepository) GetAll() ([]aggregates.Product, error) {
	var products []aggregates.Product

	postgresRepo.database.Find(&products)

	if length := len(products); length == 0 {
		return products, product.ErrorThereAreNoProducts
	}

	return products, nil
}

// GetByID finds a product by its ID
func (postgresRepo *PostgresRepository) GetByID(id uuid.UUID) (aggregates.Product, error) {
	var foundProduct aggregates.Product

	result := postgresRepo.database.Find(&foundProduct, id)

	if result.RowsAffected == 0 {
		return aggregates.Product{}, product.ErrorProductNotFound
	}

	return foundProduct, nil
}

// GetByName finds all products by their name
func (postgresRepo *PostgresRepository) GetByName(name string) ([]aggregates.Product, error) {
	products := []aggregates.Product{}

	result := postgresRepo.database.Where("name = ?", name).Find(&products)

	if result.RowsAffected == 0 {
		return []aggregates.Product{}, product.ErrorProductNotFound
	}

	return products, nil
}

// GetByNameAndPlatform finds a product by its name and platform
func (postgresRepo *PostgresRepository) GetByNameAndPlatform(name string, platform string) (aggregates.Product, error) {
	var foundProduct aggregates.Product

	result := postgresRepo.database.Where("name = ? AND platform = ?", name, platform).Find(&foundProduct)

	if result.RowsAffected == 0 {
		return aggregates.Product{}, product.ErrorProductNotFound
	}

	return foundProduct, nil
}

// Add will add a new product to the repository
func (postgresRepo *PostgresRepository) Add(productToAdd aggregates.Product) error {
	result := postgresRepo.database.Create(&productToAdd)

	if result.Error != nil {
		return product.ErrorFailedToAddProduct
	}

	return nil
}

// Update will replace an existing product information with the new product information
func (postgresRepo *PostgresRepository) Update(updatedProduct aggregates.Product) error {
	result := postgresRepo.database.Where("id = ?", updatedProduct.GetID()).Updates(&updatedProduct)

	if result.Error != nil {
		return product.ErrorUpdateProduct
	}

	return nil
}
