package main

import (
	"html/template"
	"io/fs"
	"math"
	"path/filepath"
	"time"

	uinterface "github.com/vnxcius/ecommerce-go/interface"
	"github.com/vnxcius/ecommerce-go/internal/models"
)

type templateData struct {
	Product               models.Product
	Products              []models.Product
	ProductItem           models.ProductItem
	ProductView           models.ProductView
	ProductsView          []models.ProductView
	ProductCategories     []models.ProductCategory
	ProductSize           models.ProductSize
	ProductSizes          []models.ProductSize
	ProductSizeCategory   models.ProductSizeCategory
	ProductSizeCategories []models.ProductSizeCategory
	ProductColor          models.ProductColor
	ProductColors         []models.ProductColor
	ProductImage          models.ProductImage
	ProductImages         []models.ProductImage
	ProductVariation      models.ProductVariation
	ProductVariations     []models.ProductVariation

	User      models.User
	Users     []models.User
	Address   models.Address
	Addresses []models.Address

	Form            any
	CurrentYear     int
	Version         string
	SuccessFlash    string
	IsAuthenticated bool
	CSRFToken       string
	SearchQuery     string
	Domain          string
}

func add(a, b int) int {
	if a < 0 || b < 0 {
		return 0
	}

	return a + b
}

func humanDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}

	return t.Format("02 Jan 2006 às 15:04")
}

func humanDateShort(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("02/01/2006")
}

func dateFormat(t time.Time) string {
	if t.IsZero() {
		return ""
	}

	return t.Format("2006-01-02")
}

func convertPrice(price float64) float64 {
	price = math.Round(price*100) / 100
	return price
}

func calculateDiscount(price float64, discount float64) float64 {
	if discount < 0 {
		discount = 0
	}
	discount = price * (discount / 100)
	discountedPrice := price - discount

	// Converter para 2 casas decimais
	discountedPrice = math.Round(discountedPrice*100) / 100
	return discountedPrice
}

var functions = template.FuncMap{
	"add":               add,
	"humanDate":         humanDate,
	"humanDateShort":    humanDateShort,
	"dateFormat":        dateFormat,
	"convertPrice":      convertPrice,
	"calculateDiscount": calculateDiscount,
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := fs.Glob(uinterface.Files, "html/pages/**/*.tmpl.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		patterns := []string{
			"html/pages/base.tmpl.html",
			"html/partials/**/*.tmpl.html",
			"html/partials/*.tmpl.html",
			page,
		}

		ts, err := template.New(name).Funcs(functions).ParseFS(uinterface.Files, patterns...)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
