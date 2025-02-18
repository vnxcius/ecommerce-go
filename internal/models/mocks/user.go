package mocks

import (
	"time"

	"github.com/vnxcius/ecommerce-go/internal/models"
	"gorm.io/gorm"
)

var mockUser = models.User{
	Email:     "john.doe@example.com",
	Password:  "password",
	FullName:  "John Doe",
	BirthDate: time.Now(),
	CPF:       "99888999855",
	Phone:     "99999999991",
	RoleID:    1,
}

type UserModel struct{}

func (m *UserModel) Create(user models.User) error {
	if user.Email == "dupe@example.com" {
		return gorm.ErrDuplicatedKey
	}

	if user.CPF == "11111111111" {
		return gorm.ErrDuplicatedKey
	}

	if user.Phone == "11111111111" {
		return gorm.ErrDuplicatedKey
	}

	return nil
}

func (m *UserModel) GetAll() ([]models.User, error) {
	return nil, nil
}

func (m *UserModel) GetInfo(id int) (models.User, error) {
	return models.User{}, nil
}

func (m *UserModel) Exists(id int) (bool, error) {
	return false, nil
}

func (m *UserModel) Update(id int, user models.User) error {
	return nil
}

func (m *UserModel) PermaDelete(id int) error {
	return nil
}

func (m *UserModel) ChangePassword(id int, currentPassword, password string) error {
	return nil
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	if email == "john.doe@example.com" && password == "password" {
		return 1, nil
	}

	return 0, models.ErrInvalidCredentials
}

func (m *UserModel) UpdatePhoto(id int, profilePic string) error {
	return nil
}

func (m *UserModel) IsAdmin(id int) (bool, error) {
	return false, nil
}
