package mocks

import "github.com/vnxcius/ecommerce-go/internal/models"

type ProductModel struct{}

var mockProduct = models.Product{
	Name:        "Product 1",
	Description: "Product 1 description",
	Slug:        "product-1",
	CategoryID:  1,
}

func (m *ProductModel) Get(slug string) (models.Product, error) {
	if slug == "product-1" {
		return mockProduct, nil
	}

	return models.Product{}, models.ErrNoRecords
}

func (m *ProductModel) GetAll() ([]models.Product, error) {
	return []models.Product{mockProduct}, nil
}

func (m *ProductModel) Create(name, description string, category int32) (models.Product, error) {
	return models.Product{}, nil
}

func (m *ProductModel) Update(id int, name string, description string) error {
	return nil
}

func (m *ProductModel) SoftDelete(id int) error {
	return nil
}

func (m *ProductModel) PermaDelete(id uint) error {
	return nil
}

func (m *ProductModel) Latest() ([]models.Product, error) {
	return []models.Product{mockProduct}, nil
}

func (m *ProductModel) GetFromView(slug string) (models.ProductView, error) {
	return models.ProductView{}, nil
}

func (m *ProductModel) GetFromCategory(category string) ([]models.ProductView, error) {
	return []models.ProductView{}, nil
}

func (m *ProductModel) GetFromCategoryParent(category string) ([]models.ProductView, error) {
	return []models.ProductView{}, nil
}

func (m *ProductModel) GetAllProducts() ([]models.ProductView, error) {
	return []models.ProductView{}, nil
}

func (m *ProductModel) GetFromSearch(search string) ([]models.ProductView, error) {
	return []models.ProductView{}, nil
}
