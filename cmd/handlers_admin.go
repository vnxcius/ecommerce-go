package main

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/vnxcius/ecommerce-go/internal/models"
	"github.com/vnxcius/ecommerce-go/internal/validators"
)

/*
Os handlers das páginas são divididos em
View e Post como sufixo de cada função.

View: são handlers para exibir informações
ao usuário.

Post: são handlers para processar dados
enviados pelo usuário.
*/

func (app *application) registerAdminView(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = models.UserSignupForm{}

	app.render(w, r, http.StatusOK, "register-admin.tmpl.html", data)
}

func (app *application) registerAdminPost(w http.ResponseWriter, r *http.Request) {
	var form models.UserSignupForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	form.CheckField(validators.NotBlank(form.FullName), "fullName", "Preencha o nome completo")
	form.CheckField(validators.NotBlank(form.FirstName), "firstName", "Preencha o primeiro nome")
	form.CheckField(validators.NotBlank(form.LastName), "lastName", "Preencha o sobrenome")
	form.CheckField(validators.NotBlank(form.Email), "email", "Preencha o email")
	form.CheckField(validators.NotBlank(form.CPF), "cpf", "Preencha o cpf")
	form.CheckField(validators.NotBlank(form.Phone), "phone", "Preencha o telefone")
	form.CheckField(validators.NotBlank(form.BirthDate), "birthDate", "Preencha a data de nascimento")
	form.CheckField(validators.NotBlank(form.Password), "password", "Preencha a senha")

	form.CheckField(validators.Matches(form.FullName, validators.NameRX), "fullName", "Nome completo inválido")
	form.CheckField(validators.Matches(form.FirstName, validators.NameRX), "firstName", "Nome inválido")
	form.CheckField(validators.Matches(form.LastName, validators.NameRX), "lastName", "Sobrenome inválido")
	form.CheckField(validators.Matches(form.Email, validators.EmailRX), "email", "Email inválido")
	form.CheckField(validators.Matches(form.CPF, validators.NumberRX), "cpf", "CPF inválido")
	form.CheckField(validators.Matches(form.Phone, validators.NumberRX), "phone", "Telefone inválido")

	form.CheckField(validators.MinChars(form.Password, 8), "password", "A senha deve conter ao menos 8 caracteres")
	form.CheckField(validators.IsEqual(form.Password, form.ConfirmPassword), "confirmPassword", "As senhas devem ser iguais")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form

		app.render(w, r, http.StatusUnprocessableEntity, "register-admin.tmpl.html", data)
		return
	}

	date, _ := time.Parse("2006-01-02", form.BirthDate)
	user := models.User{
		FirstName: form.FirstName,
		LastName:  form.LastName,
		FullName:  form.FullName,
		Email:     form.Email,
		Password:  form.Password,
		CPF:       form.CPF,
		Phone:     form.Phone,
		BirthDate: date,
		RoleID:    1,
	}

	app.logger.Info("Criando conta de usuário admin...", "lastName", user.LastName)
	err = app.users.Create(user)
	if err != nil {
		form.AddNonFieldError(err.Error())

		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "register-admin.tmpl.html", data)
		return
	}

	app.logger.Info("Usuário criado com sucesso!", "lastName", form.LastName)
	app.sessionManager.Put(r.Context(), "successFlash", "Usuário admin criado com sucesso!")
	http.Redirect(w, r, "/usuarios/registrar-admin", http.StatusSeeOther)
}

func (app *application) loginAdminView(w http.ResponseWriter, r *http.Request) {
	// Verifica se o usuário está autenticado. Caso esteja, redireciona para o dashboard.
	if app.isAuthenticated(r) {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}

	data := app.newTemplateData(r)
	data.Form = models.UserLoginForm{}
	app.render(w, r, http.StatusOK, "login-admin.tmpl.html", data)
}

func (app *application) loginAdminPost(w http.ResponseWriter, r *http.Request) {
	var form models.UserLoginForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.logger.Error(err.Error())
		app.serverError(w, r, err)
		return
	}

	form.CheckField(validators.NotBlank(form.Email), "email", "Preencha este campo")
	form.CheckField(validators.NotBlank(form.Password), "password", "Preencha este campo")
	form.CheckField(validators.Matches(form.Email, validators.EmailRX), "email", "Endereço de email inválido")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "login-admin.tmpl.html", data)
		return
	}

	id, err := app.users.Authenticate(form.Email, form.Password)
	if err != nil {
		form.AddNonFieldError("Email ou senha incorretos")

		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "login-admin.tmpl.html", data)
		return
	}

	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	app.sessionManager.Put(r.Context(), userIDKey, id)
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func (app *application) updateAdminPhoto(w http.ResponseWriter, r *http.Request) {
	var form models.AdminUpdateForm

	err := app.decodeMultipartPostForm(r, &form)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// Pegar arquivo
	file, _, err := r.FormFile("image")
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	defer file.Close()

	// Deletar imagem do Cloudinary
	user, _ := app.users.GetInfo(app.sessionManager.GetInt(r.Context(), userIDKey))
	err = app.deleteProfileImage(user.ProfilePic)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// Upload da imagem para o Cloudinary
	imageName, err := app.uploadProfileImage(file, "admin")
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.logger.Info("Atualizando imagem de perfil admin...", "image", imageName)
	err = app.users.UpdatePhoto(app.sessionManager.GetInt(r.Context(), userIDKey), imageName)

	if err != nil {
		app.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) usersView(w http.ResponseWriter, r *http.Request) {
	users, err := app.users.GetAll()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := app.newTemplateData(r)
	data.Users = users
	app.render(w, r, http.StatusOK, "users.tmpl.html", data)
}

func (app *application) dashboardView(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.render(w, r, http.StatusOK, "dashboard.tmpl.html", data)
}

func (app *application) allProductsView(w http.ResponseWriter, r *http.Request) {
	products, err := app.productItem.GetAllProducts(-1)

	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := app.newTemplateData(r)
	data.ProductsView = products

	app.render(w, r, http.StatusOK, "product-list.tmpl.html", data)
}

func (app *application) productView(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	productCode := params["product_code"]

	// Requisitar da view
	product, err := app.products.GetFromView(productCode)
	if err != nil {
		app.notFound(w, r)
		app.logger.Error(err.Error())
		return
	}

	data := app.newTemplateData(r)
	data.ProductView = product

	app.render(w, r, http.StatusOK, "product-edit.tmpl.html", data)
}

func (app *application) createProductView(w http.ResponseWriter, r *http.Request) {
	categories, _ := app.productCategories.GetAll()
	colors, _ := app.productColors.GetAll()
	sizes, _ := app.productSizes.GetAll()
	sizeCategories, _ := app.productSizeCategories.GetAll()

	data := app.newTemplateData(r)
	data.ProductCategories = categories
	data.ProductColors = colors
	data.ProductSizes = sizes
	data.ProductSizeCategories = sizeCategories
	data.Form = models.ProductForm{}
	app.render(w, r, http.StatusOK, "product-create.tmpl.html", data)
}

func (app *application) categoriesView(w http.ResponseWriter, r *http.Request) {
	categories, err := app.productCategories.GetAll()

	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := app.newTemplateData(r)
	data.ProductCategories = categories
	data.Form = models.CategoryForm{}
	app.render(w, r, http.StatusOK, "product-categories.tmpl.html", data)
}

func (app *application) sizesView(w http.ResponseWriter, r *http.Request) {
	sizes, _ := app.productSizes.GetAll()
	sizeCategories, err := app.productSizeCategories.GetAll()

	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := app.newTemplateData(r)
	data.ProductSizes = sizes
	data.ProductSizeCategories = sizeCategories
	data.Form = models.SizeForm{}
	app.render(w, r, http.StatusOK, "product-sizes.tmpl.html", data)
}

func (app *application) colorsView(w http.ResponseWriter, r *http.Request) {
	colors, _ := app.productColors.GetAll()

	data := app.newTemplateData(r)
	data.ProductColors = colors
	data.Form = models.ColorForm{}
	app.render(w, r, http.StatusOK, "product-colors.tmpl.html", data)
}

func (app *application) ordersView(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.render(w, r, http.StatusOK, "orders.tmpl.html", data)
}
