package models

import (
	"gorm.io/gorm"
)

type ProductItem struct {
	gorm.Model

	Price    float64 `gorm:"notnull;"`
	Discount int32   `gorm:"default:0"`
	Code     string  `gorm:"notnull;unique;type:varchar(25);index;"`
	Active   bool    `gorm:"notnull;default:false"`

	ProductID uint
	Product   Product `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE;"`

	ColorID int32
	Color   ProductColor `gorm:"foreignKey:ColorID;constraint:OnDelete:SET NULL;"`
}

type ProductItemModel struct {
	DB *gorm.DB
}

func (ProductItem) TableName() string {
	return "products.product_item"
}

func (m *ProductItemModel) Create(productItem ProductItem) (uint, error) {
	result := m.DB.Create(&productItem)

	if result.Error != nil {
		return 0, result.Error
	}

	return productItem.ID, nil
}

func (m *ProductItemModel) GetAllProducts(limit int) ([]ProductView, error) {
	var products []ProductView
	result := m.DB.
		Table("products.product_item pi").
		Order("pi.created_at DESC").
		Limit(limit).
		Select(
			"DISTINCT ON (pi.created_at) "+ // Por algum motivo Distinct() não funciona, então fiz no Select()
			"pi.id",
			"pi.code AS product_code",
			"p.name",
			"p.description",
			"p.slug",
			"pim.image",
			"pc.name AS category",
			"pc.image AS category_image",
			"pc.slug AS category_slug",
			"c.id AS color_id",
			"c.name AS color",
			"s.product_size_category_id AS size_category_id",
			"s.id AS size_id",
			"s.name AS size",
			"pv.in_stock",
			"pi.price",
			"pi.discount",
			"pi.active",
			"pi.created_at AS created_at",
			"pi.updated_at AS updated_at",
			"pi.deleted_at AS deleted_at",
		).
		Joins("JOIN products.product p ON pi.product_id = p.id").
		Joins("JOIN products.product_color c ON pi.color_id = c.id").
		Joins("JOIN products.product_category pc ON pc.id = p.category_id").
		Joins("JOIN products.product_variation pv ON pi.id = pv.product_item_id").
		Joins("JOIN products.product_size s ON pv.size_id = s.id").
		Joins("JOIN products.product_image pim ON pim.product_item_id = pi.id").
		Scan(&products)

	return products, result.Error
}
