package models

import "gorm.io/gorm"

type ProductVariation struct {
	ID uint `gorm:"primaryKey;"`

	ProductItemID uint
	ProductItem   ProductItem `gorm:"foreignKey:ProductItemID;constraint:OnDelete:CASCADE;"`

	SizeID uint
	Size   ProductSize `gorm:"foreignKey:SizeID;constraint:OnDelete:CASCADE;"`

	InStock int32
}

type ProductVariationModel struct {
	DB *gorm.DB
}

func (ProductVariation) TableName() string {
	return "products.product_variation"
}

func (m *ProductVariationModel) Create(productItemID, sizeID uint, inStock int32) error {
	var productVariation = ProductVariation{
		ProductItemID: productItemID,
		SizeID:        sizeID,
		InStock:       inStock,
	}

	return m.DB.Create(&productVariation).Error
}

func (m *ProductVariationModel) Update(id int, productItemID, sizeID uint, inStock int32) error {
	var productVariation = ProductVariation{
		ProductItemID: productItemID,
		SizeID:        sizeID,
		InStock:       inStock,
	}

	return m.DB.Model(&ProductVariation{}).Where("id = ?", id).Updates(&productVariation).Error
}

func (m *ProductVariationModel) Delete(id int) error {
	return m.DB.Unscoped().Where("id = ?", id).Delete(&ProductVariation{}).Error
}

func (m *ProductVariationModel) GetSizes(id int) ([]ProductVariation, error) {
	var variations []ProductVariation
	err := m.DB.Where("product_item_id = ?", id).Find(&variations).Error
	return variations, err
}
