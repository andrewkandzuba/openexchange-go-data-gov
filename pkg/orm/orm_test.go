package orm

import (
	"github.com/jinzhu/gorm"
	"testing"
)

type Product struct {
	gorm.Model
	Code string
	Price uint
}

func TestModel_Persist_Success(t *testing.T) {
	// Migrate the schema
	db.AutoMigrate(&Product{})

	// Create
	db.Create(&Product{Code: "L1212", Price: 1000})

	// Read
	var product Product
	db.First(&product, 1) // find product with id 1
	db.First(&product, "code = ?", "L1212") // find product with code l1212

	// Update - update product's price to 2000
	db.Model(&product).Update("Price", 2000)

	// Delete - delete product
	db.Delete(&product)
}
