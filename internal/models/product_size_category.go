package models

import "gorm.io/gorm"

type ProductSizeCategory struct {
	gorm.Model
	Name string `gorm:"type:varchar(25);notnull;unique"`
}

type ProductSizeCategoryModel struct {
	DB *gorm.DB
}

func (ProductSizeCategory) TableName() string {
	return "products.product_size_category"
}

func (m *ProductSizeCategoryModel) Create(name string) error {
	return m.DB.Create(&ProductSizeCategory{Name: name}).Error
}

func (m *ProductSizeCategoryModel) Update(id int, name string) error {
	return m.DB.Model(&ProductSizeCategory{}).Where("id = ?", id).
		Updates(&ProductSizeCategory{Name: name}).Error
}

func (m *ProductSizeCategoryModel) GetAll() ([]ProductSizeCategory, error) {
	var sizeCategories []ProductSizeCategory
	err := m.DB.Find(&sizeCategories).Error
	return sizeCategories, err
}

func (m *ProductSizeCategoryModel) PermaDelete(id int) error {
	return m.DB.Where("id = ?", id).Unscoped().
		Delete(&ProductSizeCategory{}).Error
}

func (m *ProductSizeCategoryModel) SoftDelete(id int) error {
	return m.DB.Where("id = ?", id).Delete(&ProductSizeCategory{}).Error
}
