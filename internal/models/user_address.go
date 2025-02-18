package models

import "gorm.io/gorm"

type UserAddress struct {
	ID        uint `gorm:"primarykey"`
	UserID    int  `gorm:"notnull"`
	AddressID uint `gorm:"notnull"`

	User    User    `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
	Address Address `gorm:"foreignKey:AddressID;constraint:OnDelete:CASCADE;"`
}

type UserAddressModel struct {
	DB *gorm.DB
}

type UserAddressInterface interface {
	Create(userAddress UserAddress) error
	PermaDelete(id int) error
}

func (UserAddress) TableName() string {
	return "users.user_address"
}

func (m *UserAddressModel) Create(userID int, addressID uint) error {
	return m.DB.Create(&UserAddress{UserID: userID, AddressID: addressID}).Error
}

func (m *UserAddressModel) PermaDelete(id int) error {
	return m.DB.Unscoped().Where("id = ?", id).Delete(&UserAddress{}).Error
}
