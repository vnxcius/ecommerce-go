package models

import (
	"errors"

	"gorm.io/gorm"
)

type ProductSize struct {
	gorm.Model

	Name      string `gorm:"type:varchar(10);notnull;unique"`
	SortOrder int32

	ProductSizeCategoryID int32
	ProductSizeCategory   ProductSizeCategory `gorm:"foreignKey:ProductSizeCategoryID;constraint:OnDelete:SET NULL;"`
}

type ProductSizeModel struct {
	DB *gorm.DB
}

func (ProductSize) TableName() string {
	return "products.product_size"
}

func (m *ProductSizeModel) Create(name string, sortOrder int32, productSizeCategoryID int32) error {
	var size = ProductSize{
		Name:                  name,
		SortOrder:             sortOrder,
		ProductSizeCategoryID: productSizeCategoryID,
	}

	// Verificar se o tamanho existe
	result := m.DB.Where("name = ?", name).First(&size)
	if result.RowsAffected > 0 {
		return errors.New("tamanho já existente")
	}

	return m.DB.Create(&size).Error
}

func (m *ProductSizeModel) Update(id int, name string, sortOrder int32) error {
	return m.DB.Model(&ProductSize{}).Where("id = ?", id).
		Updates(&ProductSize{Name: name, SortOrder: sortOrder}).Error
}

func (m *ProductSizeModel) GetAll() ([]ProductSize, error) {
	var sizes []ProductSize
	err := m.DB.Joins("ProductSizeCategory").Order("sort_order ASC").Find(&sizes).Error
	return sizes, err
}

func (m *ProductSizeModel) SoftDelete(id int) error {
	return m.DB.Where("id = ?", id).Delete(&ProductSize{}).Error
}

func (m *ProductSizeModel) GetAllByProduct(id int) ([]ProductSize, error) {
	var sizes []ProductSize
	err := m.DB.Where("product_size_category_id = ?", id).Order("sort_order ASC").Find(&sizes).Error
	return sizes, err
}
