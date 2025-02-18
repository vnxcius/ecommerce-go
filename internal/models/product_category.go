package models

import (
	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

// Possui o campo ParentCategory como um self-join para casos onde uma categoria possui categorias filhas.
type ProductCategory struct {
	gorm.Model

	Name        string `gorm:"type:varchar(50);notnull;"`
	Image       string `gorm:"type:varchar(255);"`
	Description string `gorm:"type:varchar(255);"`
	Slug        string `gorm:"type:varchar(255);notnull;unique;index;"`

	ParentCategoryID int32            `gorm:"default:null"`
	ParentCategory   *ProductCategory `gorm:"foreignKey:ParentCategoryID;constraint:OnDelete:SET NULL;"`
}

type ProductCategoryModel struct {
	*gorm.DB
}

func (ProductCategory) TableName() string {
	return "products.product_category"
}

func (m *ProductCategoryModel) GetAll() ([]ProductCategory, error) {
	var categories []ProductCategory

	result := m.DB.Joins("ParentCategory").Find(&categories)

	if result.Error != nil {
		return nil, result.Error
	}

	return categories, nil
}

func (m *ProductCategoryModel) Create(categoryName string, categoryDescription string, categoryImage string, parentCategoryID int) error {
	var category = ProductCategory{
		Name:             categoryName,
		Description:      categoryDescription,
		Image:            categoryImage,
		Slug:             slug.Make(categoryName),
		ParentCategoryID: int32(parentCategoryID),
	}

	return m.DB.Create(&category).Error
}

func (m *ProductCategoryModel) PermaDelete(id int) error {
	return m.DB.Unscoped().Where("id = ?", id).Delete(&ProductCategory{}).Error
}

func (m *ProductCategoryModel) SoftDelete(id int) error {
	return m.DB.Where("id = ?", id).Delete(&ProductCategory{}).Error
}
