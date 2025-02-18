package models

import (
	"gorm.io/gorm"
)

type ProductColor struct {
	gorm.Model

	Name string `gorm:"type:varchar(20);notnull;unique"`
	Hex  string `gorm:"type:varchar(7);notnull;unique;"`
}

type ProductColorModel struct {
	DB *gorm.DB
}

func (ProductColor) TableName() string {
	return "products.product_color"
}

func (m *ProductColorModel) GetAll() ([]ProductColor, error) {
	var colors []ProductColor
	err := m.DB.Find(&colors).Error
	return colors, err
}

func (m *ProductColorModel) Create(name string, hex string) error {
	return m.DB.Create(&ProductColor{Name: name, Hex: hex}).Error
}

func (m *ProductColorModel) Update(id int, name string) error {
	var color = ProductColor{
		Name: name,
	}

	return m.DB.Model(&ProductColor{}).Where("id = ?", id).Updates(&color).Error
}

func (m *ProductColorModel) Delete(id int) error {
	return m.DB.Unscoped().Where("id = ?", id).Delete(&ProductColor{}).Error
}
