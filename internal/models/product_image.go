package models

import (
	"gorm.io/gorm"
)

type ProductImage struct {
	gorm.Model

	Image string `gorm:"type:varchar(255);notnull;unique"`

	ProductItemID uint
	Product       ProductItem `gorm:"foreignKey:ProductItemID;constraint:OnDelete:CASCADE;"`
}

type ProductImageModel struct {
	DB *gorm.DB
}

func (ProductImage) TableName() string {
	return "products.product_image"
}

func (m *ProductImageModel) Create(productImage string, productItemID uint) error {
	var image = ProductImage{
		Image:         productImage,
		ProductItemID: productItemID,
	}

	return m.DB.Create(&image).Error
}

func (m *ProductImageModel) Delete(id int) error {
	return m.DB.Unscoped().Where("id = ?", id).Delete(&ProductImage{}).Error
}

func (m *ProductImageModel) GetAll() ([]ProductImage, error) {
	var images []ProductImage
	err := m.DB.Find(&images).Error
	return images, err
}

func (m *ProductImageModel) GetAllByProduct(id int) ([]ProductImage, error) {
	var images []ProductImage
	err := m.DB.Where("product_item_id = ?", id).Find(&images).Error
	return images, err
}
