package models

type UserRole struct {
	ID uint `gorm:"primarykey"`

	Name string `gorm:"type:varchar(10);notnull;unique"`
}

func (UserRole) TableName() string {
	return "users.user_role"
}
