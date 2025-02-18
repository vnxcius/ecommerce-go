package models

import (
	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model

	Name        string `gorm:"type:varchar(255);notnull;"`
	Description string `gorm:"type:text;notnull"`
	Slug        string `gorm:"type:varchar(255);notnull;unique;index;"`

	CategoryID int32
	Category   ProductCategory `gorm:"foreignKey:CategoryID;constraint:OnDelete:SET NULL;"`
}

type ProductInterface interface {
	Create(name, description string, category int32) (Product, error)
	Update(id int, name string, description string) error
	Get(productCode string) (Product, error)
	GetFromView(slug string) (ProductView, error)
	GetFromCategory(category string) ([]ProductView, error)
	GetFromCategoryParent(category string) ([]ProductView, error)
	GetFromSearch(search string) ([]ProductView, error)
	GetAll() ([]Product, error)
	Latest() ([]Product, error)
	PermaDelete(id uint) error
	SoftDelete(id int) error
}

type ProductModel struct {
	DB *gorm.DB
}

func (Product) TableName() string {
	return "products.product"
}

func (m *ProductModel) Create(name, description string, category int32) (Product, error) {
	var product = Product{
		Name:        name,
		Description: description,
		Slug:        slug.Make(name),
		CategoryID:  category,
	}

	result := m.DB.Where("Name = ?", name).First(&product)

	if result.RowsAffected > 0 {
		return Product{}, result.Error
	}

	result = m.DB.Create(&product)

	if result.Error != nil {
		return Product{}, result.Error
	}

	return product, nil
}

func (m *ProductModel) Update(id int, name string, description string) error {
	return nil
}

func (m *ProductModel) Get(productCode string) (Product, error) {
	return Product{}, nil
}

func (m *ProductModel) GetAll() ([]Product, error) {
	var products []Product

	result := m.DB.Find(&products)
	if result.Error != nil {
		return nil, result.Error
	}
	return products, nil
}

func (p *ProductModel) Latest() ([]Product, error) {
	return nil, nil
}

func (m *ProductModel) PermaDelete(id uint) error {
	var product Product
	return m.DB.Unscoped().Where("id = ?", id).Delete(&product).Error
}

func (m *ProductModel) SoftDelete(id int) error {
	var product Product
	return m.DB.Where("id = ?", id).Delete(&product).Error
}
