package models

import (
	"mime/multipart"

	"github.com/vnxcius/ecommerce-go/internal/validators"
)

type ProductForm struct {
	Image       multipart.FileHeader
	Images      []multipart.File
	Name        string
	Description string
	Category    int32
	Size        int32
	Color       int32
	Price       float64
	Discount    int32
	InStock     int32
	Active      bool

	validators.Validator
}

type CategoryForm struct {
	Name        string
	Description string
	Parent      string

	validators.Validator
}

type SizeForm struct {
	Name      string
	NameSize  string
	SortOrder string
	Parent    string

	validators.Validator
}

type ColorForm struct {
	Name string
	Hex  string

	validators.Validator
}

type UserSignupForm struct {
	FirstName       string
	LastName        string
	FullName        string
	CPF             string
	Phone           string
	Email           string
	Password        string
	ConfirmPassword string
	BirthDate       string

	validators.Validator
}

type UserLoginForm struct {
	Email    string
	Password string

	validators.Validator
}

type ChangePasswordForm struct {
	CurrentPassword string
	NewPassword     string
	ConfirmPassword string

	validators.Validator
}

type DeleteAccountForm struct {
	validators.Validator
}

type AdminUpdateForm struct {
	FirstName       string
	LastName        string
	Email           string
	Password        string
	ConfirmPassword string
	ProfilePic      string

	validators.Validator
}

type ShippingForm struct {
	CEP string

	validators.Validator
}

type SearchForm struct {
	SearchQuery string

	validators.Validator
}

type AddressForm struct {
	FullName   string
	CEP        string
	Street     string
	Number     string
	Complement string
	District   string
	City       string
	UF         string
	Reference  string

	validators.Validator
}
