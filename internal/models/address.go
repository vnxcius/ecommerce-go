package models

import "gorm.io/gorm"

type Address struct {
	gorm.Model

	HashID     string `gorm:"type:varchar(101);notnull;unique;index"`
	Name       string `gorm:"type:varchar(50);notnull"`
	Cep        string `gorm:"type:varchar(50);notnull"`
	Street     string `gorm:"type:varchar(50);notnull"`
	Number     string `gorm:"type:varchar(50);notnull"`
	Complement string `gorm:"type:varchar(50);"`
	District   string `gorm:"type:varchar(50);notnull"`
	City       string `gorm:"type:varchar(50);notnull"`
	UF         string `gorm:"type:varchar(2);notnull"`
	Reference  string `gorm:"type:varchar(50);"`
}

type AddressModel struct {
	DB *gorm.DB
}

func (Address) TableName() string {
	return "users.address"
}

func (m *AddressModel) Create(Address Address) (uint, error) {
	result := m.DB.Create(&Address)

	if result.Error != nil {
		return 0, result.Error
	}

	return Address.ID, nil
}

func (m *AddressModel) Update(hashID string, address Address) error {
	return m.DB.Model(&Address{}).Where("hash_id = ?", hashID).Updates(&address).Error
}

func (m *AddressModel) PermaDeleteByHashID(id string) error {
	return m.DB.Unscoped().Where("hash_id = ?", id).Delete(&Address{}).Error
}

func (m *AddressModel) AddHashID(addressID uint, hashID string) error {
	return m.DB.Model(&Address{}).Where("id = ?", addressID).Update("hash_id", hashID).Error
}

func (m *AddressModel) GetAll(userID int) ([]Address, error) {
	var addresses []Address
	err := m.DB.
		Table("users.user_address").
		Joins("JOIN users.address ad ON ad.id  = user_address.address_id").
		Where("user_address.user_id = ?", userID).
		Scan(&addresses).Error
	return addresses, err
}

func (m *AddressModel) GetByHashID(hashID string) (Address, error) {
	var address Address
	err := m.DB.Where("hash_id = ?", hashID).First(&address).Error
	return address, err
}
