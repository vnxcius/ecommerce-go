package models

import (
	"html/template"
	"time"
)

type ProductView struct {
	ID                  uint          `gorm:"primarykey;column:id"`
	ProductCode         string        `gorm:"column:product_code"`
	Name                string        `gorm:"column:name"`
	Description         template.HTML `gorm:"column:description"`
	Slug                string        `gorm:"column:slug"`
	Image               string        `gorm:"column:image"`
	CategoryImage       string        `gorm:"column:category_image"`
	Category            string        `gorm:"column:category"`
	CategorySlug        string        `gorm:"column:category_slug"`
	CategoryParent      string        `gorm:"column:category_parent"`
	CategoryParentImage string        `gorm:"column:category_parent_image"`
	CategoryParentSlug  string        `gorm:"column:category_parent_slug"`
	ColorID             int32         `gorm:"column:color_id"`
	Color               string        `gorm:"column:color"`
	SizeCategoryID      int32         `gorm:"column:size_category_id"`
	SizeID              int32         `gorm:"column:size_id"`
	Size                string        `gorm:"column:size"`
	InStock             string        `gorm:"column:in_stock"`
	Price               float64       `gorm:"column:price"`
	Discount            float64       `gorm:"column:discount"`
	Active              bool          `gorm:"column:active"`
	CreatedAt           time.Time     `gorm:"column:created_at"`
	UpdatedAt           time.Time     `gorm:"column:updated_at"`
	DeletedAt           time.Time     `gorm:"column:deleted_at"`
}

func (ProductView) TableName() string {
	return "products.products_view"
}

func (m *ProductModel) GetAllFromView() ([]ProductView, error) {
	var products []ProductView

	result := m.DB.Find(&products)
	if result.Error != nil {
		return nil, result.Error
	}
	return products, nil
}

func (m *ProductModel) GetFromView(slug string) (ProductView, error) {
	var productView ProductView

	result := m.DB.Where("slug = ?", slug).First(&productView)

	if result.Error != nil {
		return ProductView{}, result.Error
	}

	return productView, nil
}

func (m *ProductModel) GetFromCategory(slug string) ([]ProductView, error) {
	var products []ProductView
	result := m.DB.Where("category_slug = ?", slug).Find(&products)

	return products, result.Error
}

func (m *ProductModel) GetFromCategoryParent(slug string) ([]ProductView, error) {
	var products []ProductView
	result := m.DB.Where("category_parent_slug = ?", slug).Find(&products)

	return products, result.Error
}

func (m *ProductModel) GetFromSearch(search string) ([]ProductView, error) {
	var products []ProductView
	result := m.DB.Where("SIMILARITY(name, ?) >= 0.1", search).Order("created_at DESC").Find(&products)

	return products, result.Error
}
