package models

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	ProfilePic string    `gorm:"type:varchar(255);notnull;default:users/default_profile_pic.png"`
	FirstName  string    `gorm:"type:varchar(25);"`
	LastName   string    `gorm:"type:varchar(25);"`
	FullName   string    `gorm:"type:varchar(50);notnull"`
	Email      string    `gorm:"type:varchar(50);unique;notnull"`
	Password   string    `gorm:"type:varchar(255);notnull"`
	BirthDate  time.Time `gorm:"type:date;notnull"`
	CPF        string    `gorm:"type:varchar(11);unique;notnull"`
	Phone      string    `gorm:"type:varchar(11);unique;notnull"`
	RoleID     int32     `gorm:"notnull;default:2"`
	Role       UserRole  `gorm:"foreignKey:RoleID;constraint:OnDelete:SET NULL;"`
}

type UserModel struct {
	DB *gorm.DB
}

type UserModelInterface interface {
	GetAll() ([]User, error)
	GetInfo(id int) (User, error)
	Exists(id int) (bool, error)
	Create(user User) error
	PermaDelete(id int) error
	Update(id int, user User) error
	ChangePassword(id int, currentPassword, password string) error
	Authenticate(email, password string) (int, error)
	UpdatePhoto(id int, profilePic string) error
	IsAdmin(id int) (bool, error)
}

func (User) TableName() string {
	return "users.user"
}

func (m *UserModel) GetAll() ([]User, error) {
	var users []User
	err := m.DB.Order("id ASC").Find(&users).Error
	return users, err
}

func (m *UserModel) GetInfo(id int) (User, error) {
	var user User
	err := m.DB.Where("id = ?", id).First(&user).Error
	return user, err
}

func (m *UserModel) Exists(id int) (bool, error) {
	var user User
	err := m.DB.Where("id = ?", id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		} else {
			return false, err
		}
	}

	// Usuário existe
	return true, nil
}

func (m *UserModel) Create(user User) error {
	// Verificar se o email existe
	err := m.DB.Where("email = ?", user.Email).First(&user).Error
	if err == nil {
		return errors.New("este email já foi cadastrado")
	}

	// Verificar se o cpf existe
	err = m.DB.Where("cpf = ?", user.CPF).First(&user).Error
	if err == nil {
		return errors.New("este cpf já foi cadastrado")
	}

	// Verificar se o telefone existe
	err = m.DB.Where("phone = ?", user.Phone).First(&user).Error
	if err == nil {
		return errors.New("este telefone já foi cadastrado")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)

	return m.DB.Create(&user).Error
}

func (m *UserModel) Update(id int, user User) error {
	return m.DB.Model(&user).Where("id = ?", id).Updates(&user).Error
}

func (m *UserModel) PermaDelete(id int) error {
	return m.DB.Unscoped().Where("id = ?", id).Delete(&User{}).Error
}

func (m *UserModel) ChangePassword(id int, currentPassword, password string) error {
	// Verificar se a senha é igual a atual
	if currentPassword == password {
		return errors.New("a senha atual deve ser diferente da nova senha")
	}

	var user User
	err := m.DB.Where("id = ?", id).First(&user).Error
	if err != nil {
		return err
	}

	// Validar senha atual
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(currentPassword))
	if err != nil {
		// retorna invalid credentials se a senha estiver errada
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return ErrInvalidCredentials
		}

		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	return m.DB.Model(&User{}).Where("id = ?", id).Update("password", string(hashedPassword)).Error
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	var user User

	err := m.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return 0, err
	}

	// validar hash
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		// retorna invalid credentials se a senha estiver errada
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		}

		return 0, err
	}

	return int(user.ID), nil
}

func (m *UserModel) UpdatePhoto(id int, profilePic string) error {
	return m.DB.Where("id = ?", id).Updates(&User{ProfilePic: profilePic}).Error
}

func (m *UserModel) IsAdmin(id int) (bool, error) {
	var user User
	err := m.DB.Where("id = ?", id).First(&user).Error
	if err != nil {
		return false, err
	}

	if user.RoleID == 1 {
		return true, nil
	}

	return false, nil
}
