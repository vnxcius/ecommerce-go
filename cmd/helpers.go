package main

import (
	"bytes"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/google/uuid"
	"github.com/gorilla/schema"
	"github.com/justinas/nosurf"
	"github.com/speps/go-hashids/v2"
	"github.com/vnxcius/ecommerce-go/internal/models"
)

func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
	)

	app.logger.Error(err.Error(), "method", method, "uri", uri)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) notFound(w http.ResponseWriter, r *http.Request) {
	app.logger.Error("Page not found", "method", r.Method, "uri", r.RequestURI)
	app.render(w, nil, http.StatusNotFound, "404.tmpl.html", templateData{})
}

func (app *application) render(w http.ResponseWriter, r *http.Request, status int, page string, data templateData) {
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(w, r, err)
		return
	}

	buf := new(bytes.Buffer)

	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	w.WriteHeader(status)

	buf.WriteTo(w)
}

func (app *application) newTemplateData(r *http.Request) templateData {
	categories, _ := app.productCategories.GetAll()

	return templateData{
		CurrentYear:       time.Now().Year(),
		Version:           "v0.0.506",
		SuccessFlash:      app.sessionManager.PopString(r.Context(), "successFlash"),
		IsAuthenticated:   app.isAuthenticated(r),
		CSRFToken:         nosurf.Token(r),
		User:              app.getUserInfo(r),
		ProductCategories: categories,
		SearchQuery:       r.FormValue("q"),
		Domain:            app.domain,
	}
}

func (app *application) decodePostForm(r *http.Request, dst any) error {
	err := r.ParseForm()
	if err != nil {
		app.logger.Info("Error parsing form")
		return err
	}

	err = app.formDecoder.Decode(dst, r.PostForm)
	if err != nil {
		app.logger.Info("Error decoding form")
		var invalidDecoderError *schema.ConversionError

		if errors.As(err, &invalidDecoderError) {
			panic(err)
		}

		// Retornar outros erros.
		return err
	}

	return nil
}

// dst é o alvo destino que queremos decodificar os dados em.
func (app *application) decodeMultipartPostForm(r *http.Request, dst any) error {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		return err
	}

	err = app.formDecoder.Decode(dst, r.PostForm)
	if err != nil {
		var invalidDecoderError *schema.ConversionError

		if errors.As(err, &invalidDecoderError) {
			panic(err)
		}

		// Retornar outros erros.
		return err
	}
	return nil
}

func (app *application) uploadImage(file multipart.File, from string) (imageName string, err error) {
	// Gerar nome da imagem
	imageName = from + "-" + uuid.New().String()

	app.logger.Info("Upando imagem para o Cloudinary", "image", imageName)
	res, err := app.cld.Upload.Upload(app.ctx, file, uploader.UploadParams{
		PublicID:       imageName,
		AssetFolder:    "products",
		AllowedFormats: []string{"jpg", "png", "jpeg", "webp"},
	})

	if err != nil {
		return "", err
	}

	if res.Error.Message != "" {
		return "", errors.New(res.Error.Message)
	}

	return imageName, nil
}

// func (app *application) uploadMultipleImages(files []*multipart.FileHeader, productID uint, productItemID uint) error {
// 	// Para cada imagem adicional, upar para o cloudinary
// 	// e então adicionar ao banco de dados
// 	for _, image := range files {
// 		file, err := image.Open()
// 		if err != nil {
// 			return err
// 		}
// 		defer file.Close()

// 		// Upar para o cloudinary
// 		app.logger.Info("Upando imagem adicional de produto para o cloudinary...", "productID", productID)
// 		imageName, err := app.uploadImage(file, "product")
// 		if err != nil {
// 			return err
// 		}

// 		// Adicionar imagem ao banco de dados
// 		app.logger.Info("Salvando imagem adicional no banco de dados...", "productID", productID)
// 		err = app.productImage.Create(imageName, productItemID)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

func (app *application) uploadProfileImage(file multipart.File, from string) (imageName string, err error) {
	// Gerar nome da imagem
	imageName = from + "-" + uuid.New().String()

	// Upar imagem para o Cloudinary
	app.logger.Info("Upando imagem para o Cloudinary", "image", imageName)
	res, _ := app.cld.Upload.Upload(app.ctx, file, uploader.UploadParams{
		PublicID:       imageName,
		AssetFolder:    "users",
		AllowedFormats: []string{"jpg", "png", "jpeg", "webp"},
	})

	if res.Error.Message != "" {
		return "", errors.New(res.Error.Message)
	}

	return imageName, nil
}

func (app *application) deleteProfileImage(imageName string) error {
	app.logger.Info("Deletando imagem do Cloudinary", "image", imageName)
	// Deletar imagem do Cloudinary
	res, _ := app.cld.Upload.Destroy(app.ctx, uploader.DestroyParams{PublicID: imageName})

	if res.Result != "ok" {
		return errors.New(res.Error.Message)
	}

	return nil
}

func (app *application) isAuthenticated(r *http.Request) bool {
	// Checa o contexto e verifica se o user está autenticado.
	isAuthenticated, ok := r.Context().Value(isAuthenticatedUserKey).(bool)
	if !ok {
		app.logger.Error("context value for " + string(isAuthenticatedUserKey) + " not found in request context")
		return false
	}
	return isAuthenticated
}

func (app *application) getUserInfo(r *http.Request) models.User {
	userID := app.sessionManager.GetInt(r.Context(), userIDKey)
	if userID != 0 {
		user, err := app.users.GetInfo(userID)
		if err != nil {
			app.logger.Error("failed to fetch user from session", "error", err)
		}

		return user
	}

	return models.User{}
}

func (app *application) hideEmail(email string) string {
	// Dividir o email em duas partes: usuário e domínio
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return email
	}

	user := parts[0]
	domain := parts[1]

	// Determinar o número de caracteres a serem ocultados
	numAsterisks := len(user) - 1
	if numAsterisks < 1 {
		numAsterisks = 1
	}

	// Construir o novo usuário com asteriscos
	hiddenUser := string(user[0]) + strings.Repeat("*", numAsterisks)
	return hiddenUser + "@" + domain
}

func (app *application) generateHashID(id uint, length int) (string, error) {
	hd := hashids.NewData()
	hd.Salt = os.Getenv("HASH_SALT")
	hd.MinLength = length

	hash, err := hashids.NewWithData(hd)
	if err != nil {
		return "", err
	}

	result, err := hash.Encode([]int{
		int(id),
	})
	if err != nil {
		app.logger.Error("failed to encode hashids", "error", err)
	}

	return result, nil
}
