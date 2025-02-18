package main

import (
	"context"
	"flag"
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/gorilla/schema"
	"github.com/joho/godotenv"
	"github.com/vnxcius/ecommerce-go/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Define uma struct que representa a dependências da aplicação
type application struct {
	domain          string
	logger          *slog.Logger
	templateCache   map[string]*template.Template
	formDecoder     *schema.Decoder
	sessionManager  *scs.SessionManager
	MaintenanceMode bool

	users       models.UserModelInterface
	userAddress *models.UserAddressModel
	addresses   *models.AddressModel

	products              models.ProductInterface
	productItem           *models.ProductItemModel
	productImage          *models.ProductImageModel
	productCategories     *models.ProductCategoryModel
	productSizes          *models.ProductSizeModel
	productSizeCategories *models.ProductSizeCategoryModel
	productColors         *models.ProductColorModel
	productVariation      *models.ProductVariationModel

	// Cloudinary
	cld *cloudinary.Cloudinary
	ctx context.Context
}

func main() {
	// Structered logging
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Obter o caminho do arquivo .env a partir da variável de ambiente
	envFilePath := os.Getenv("ENV_FILE_PATH")
	if envFilePath == "" {
		envFilePath = "../.env" // Caminho padrão para desenvolvimento local
	}

	// Carregar as variáveis de ambiente
	err := godotenv.Load(envFilePath)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	cloudinaryUrl := os.Getenv("CLOUDINARY_URL")
	dsn := os.Getenv("DSN")
	domain := os.Getenv("DOMAIN")

	// Define a flag port e o valor padrão
	port := flag.String("port", "80", "Server port.")
	flag.Parse()

	// Conectar-se ao banco de dados
	logger.Info("Trying to connect to database...")
	db, err := openDB(dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	logger.Info("Database connection successful!")

	// Criar sqlDB
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	// Conectar-se ao cloudinary
	logger.Info("Trying to connect to Cloudinary API...")
	cld, ctx := credentialsCld(cloudinaryUrl)
	logger.Info("Cloudinary connection successful!")

	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// Iniciar instância de form decoder
	formDecoder := schema.NewDecoder()
	formDecoder.IgnoreUnknownKeys(true) // Ignorar campos dos formulários para evitar erros

	sessionManager := scs.New()
	sessionManager.Store = postgresstore.New(sqlDB)
	sessionManager.Lifetime = 24 * 30 * time.Hour // 30 dias
	sessionManager.Cookie.Secure = true

	// Iniciar nova instância da aplicação contendo as dependências necessárias
	app := &application{
		// config
		domain:          domain,
		logger:          logger,
		templateCache:   templateCache,
		formDecoder:     formDecoder,
		sessionManager:  sessionManager,
		MaintenanceMode: false,

		// cloudinary
		cld: cld,
		ctx: ctx,

		// models
		users:       &models.UserModel{DB: db},
		userAddress: &models.UserAddressModel{DB: db},
		addresses:   &models.AddressModel{DB: db},

		products:              &models.ProductModel{DB: db},
		productItem:           &models.ProductItemModel{DB: db},
		productImage:          &models.ProductImageModel{DB: db},
		productCategories:     &models.ProductCategoryModel{DB: db},
		productSizes:          &models.ProductSizeModel{DB: db},
		productSizeCategories: &models.ProductSizeCategoryModel{DB: db},
		productColors:         &models.ProductColorModel{DB: db},
		productVariation:      &models.ProductVariationModel{DB: db},
	}

	srv := &http.Server{
		Addr:         ":" + *port,
		Handler:      app.routes(),
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
		IdleTimeout:  time.Minute,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Info("Server running", "port", srv.Addr)
	err = srv.ListenAndServe()
	logger.Error(err.Error())
	os.Exit(1)
}

func credentialsCld(cldURL string) (*cloudinary.Cloudinary, context.Context) {
	cld, err := cloudinary.NewFromURL(cldURL)

	if err != nil {
		log.Fatal("ERROR: ", err)
		os.Exit(1)
	}

	cld.Config.URL.Secure = true
	ctx := context.Background()
	return cld, ctx
}

func openDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{TranslateError: true})
	if err != nil {
		return nil, err
	}

	db.Exec("CREATE EXTENSION pg_trgm;") // Criar extensão pg_trgm

	err = db.AutoMigrate(
		&models.Session{},
		&models.User{},
		&models.UserAddress{},
		&models.Address{},

		&models.ProductCategory{},
		&models.ProductColor{},
		&models.ProductImage{},
		&models.ProductItem{},
		&models.ProductSizeCategory{},
		&models.ProductSize{},
		&models.ProductVariation{},
		&models.Product{},
	)
	if err != nil {
		return nil, err
	}

	// Criar cargos padrão
	err = createRoles(db)
	if err != nil {
		return nil, err
	}

	// Criar conta de admin padrão
	err = createDefaultAdmin(db)
	if err != nil {
		return nil, err
	}

	query := db.
		Table("products.product_item pi").
		Order("pi.created_at DESC").
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
		Joins("JOIN products.product_image pim ON pim.product_item_id = pi.id")

	err = db.Migrator().CreateView("products.products_view", gorm.ViewOption{Query: query, Replace: true})

	if err != nil {
		return nil, err
	}

	return db, nil
}

func createDefaultAdmin(db *gorm.DB) error {
	var userCount int64
	if err := db.Model(&models.User{}).Count(&userCount).Error; err != nil {
		return err
	}

	if userCount > 0 {
		return nil
	}

	email := os.Getenv("DEFAULT_ADMIN_EMAIL")
	password := os.Getenv("DEFAULT_ADMIN_PASSWORD")

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	user := &models.User{
		FirstName: "Lorem",
		LastName:  "Admin",
		FullName:  "Lorem Admin",
		Email:     email,
		Password:  string(hashedPassword),
		BirthDate: time.Now(),
		Phone:     "00000000000",
		CPF:       "00000000000",
		RoleID:    1,
	}

	err = db.Create(&user).Error
	return err
}

func createRoles(db *gorm.DB) error {
	var roleCount int64

	if err := db.Model(&models.UserRole{}).Count(&roleCount).Error; err != nil {
		return err
	}

	if roleCount > 0 {
		return nil
	}

	roles := []models.UserRole{
		{Name: "admin"},
		{Name: "user"},
	}

	err := db.Create(&roles).Error
	if err != nil {
		return err
	}

	return nil
}
