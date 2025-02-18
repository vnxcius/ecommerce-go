package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/vnxcius/ecommerce-go/internal/models"
	"github.com/vnxcius/ecommerce-go/internal/validators"
	"gorm.io/gorm"
)

/*
Os handlers das páginas são divididos em
View e Post como sufixo de cada função.

View: são handlers para exibir informações
ao usuário.

Post: são handlers para processar dados
enviados pelo usuário.
*/

func (app *application) registerView(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = models.UserSignupForm{}
	app.render(w, r, http.StatusOK, "register.tmpl.html", data)
}

func (app *application) registerPost(w http.ResponseWriter, r *http.Request) {
	var form models.UserSignupForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	form.CheckField(validators.NotBlank(form.FullName), "fullName", "Preencha o nome corretamente")
	form.CheckField(validators.NotBlank(form.BirthDate), "birthDate", "Preencha a data de nascimento")
	form.CheckField(validators.NotBlank(form.CPF), "cpf", "Preencha o cpf")
	form.CheckField(validators.NotBlank(form.Phone), "phone", "Preencha o telefone")
	form.CheckField(validators.NotBlank(form.Email), "email", "Preencha o email")
	form.CheckField(validators.NotBlank(form.Password), "password", "Preencha a senha")

	form.CheckField(validators.Matches(form.FullName, validators.NameRX), "fullName", "Nome inválido")
	form.CheckField(validators.Matches(form.CPF, validators.NumberRX), "cpf", "CPF inválido")
	form.CheckField(validators.Matches(form.Phone, validators.NumberRX), "phone", "Telefone inválido")
	form.CheckField(validators.Matches(form.Email, validators.EmailRX), "email", "Email inválido")

	form.CheckField(validators.MinChars(form.Password, 8), "password", "A senha deve conter ao menos 8 caracteres")
	form.CheckField(validators.MaxChars(form.Password, 50), "password", "A senha deve conter no máximo 50 caracteres")
	form.CheckField(validators.MinChars(form.ConfirmPassword, 8), "confirmPassword", "A senha deve conter ao menos 8 caracteres")
	form.CheckField(validators.MaxChars(form.ConfirmPassword, 50), "confirmPassword", "A senha deve conter no máximo 50 caracteres")

	// Verificar se as senhas coincidem
	form.CheckField(validators.IsEqual(form.Password, form.ConfirmPassword), "confirmPassword", "As senhas devem ser iguais")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "register.tmpl.html", data)
		return
	}

	date, _ := time.Parse("2006-01-02", form.BirthDate)
	user := models.User{
		FullName:  form.FullName,
		BirthDate: date,
		CPF:       form.CPF,
		Phone:     form.Phone,
		Email:     form.Email,
		Password:  form.Password,
		RoleID:    2,
	}

	email := app.hideEmail(user.Email)

	app.logger.Info("Criando conta de usuário...", "email", email)
	err = app.users.Create(user)
	if err != nil {
		form.AddNonFieldError(err.Error())

		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "register.tmpl.html", data)
		return
	}

	app.logger.Info("Conta de usuário criada com sucesso!", "email", email)
	app.sessionManager.Put(r.Context(), "successFlash", "Conta criada com sucesso. Por favor, realize o login.")
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (app *application) loginView(w http.ResponseWriter, r *http.Request) {
	if app.isAuthenticated(r) {
		http.Redirect(w, r, "/meu-perfil", http.StatusSeeOther)
		return
	}

	data := app.newTemplateData(r)
	data.Form = models.UserLoginForm{}
	app.render(w, r, http.StatusOK, "login.tmpl.html", data)
}

func (app *application) loginPost(w http.ResponseWriter, r *http.Request) {
	var form models.UserLoginForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	form.CheckField(validators.NotBlank(form.Email), "email", "Preencha este campo")
	form.CheckField(validators.NotBlank(form.Password), "password", "Preencha este campo")

	form.CheckField(validators.Matches(form.Email, validators.EmailRX), "email", "Endereço de email inválido")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "login.tmpl.html", data)
		return
	}

	id, err := app.users.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Email ou senha incorretos")
		}

		if errors.Is(err, gorm.ErrRecordNotFound) {
			form.AddNonFieldError("Email não cadastrado")
		}

		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "login.tmpl.html", data)
		return
	}

	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	app.sessionManager.Put(r.Context(), userIDKey, id)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) homeView(w http.ResponseWriter, r *http.Request) {
	products, err := app.productItem.GetAllProducts(10)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := app.newTemplateData(r)
	data.ProductsView = products
	app.render(w, r, http.StatusOK, "home.tmpl.html", data)
}

func (app *application) searchView(w http.ResponseWriter, r *http.Request) {
	query := mux.Vars(r)["q"]
	products, err := app.products.GetFromSearch(query)

	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := app.newTemplateData(r)
	data.Form = models.SearchForm{SearchQuery: query}
	data.ProductsView = products
	app.render(w, r, http.StatusOK, "search.tmpl.html", data)
}

func (app *application) profileView(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = models.UserSignupForm{}
	app.render(w, r, http.StatusOK, "profile.tmpl.html", data)
}

func (app *application) profilePost(w http.ResponseWriter, r *http.Request) {
	var form models.UserSignupForm

	if err := app.decodePostForm(r, &form); err != nil {
		app.serverError(w, r, err)
		return
	}

	form.CheckField(validators.NotBlank(form.FullName), "fullName", "Preencha o nome corretamente")
	form.CheckField(validators.NotBlank(form.FirstName), "firstName", "Preencha o primeiro nome")
	form.CheckField(validators.NotBlank(form.LastName), "lastName", "Preencha o sobrenome")
	form.CheckField(validators.NotBlank(form.Email), "email", "Preencha o email")
	form.CheckField(validators.NotBlank(form.Phone), "phone", "Preencha o telefone")

	form.CheckField(validators.Matches(form.FullName, validators.NameRX), "fullName", "Nome inválido")
	form.CheckField(validators.Matches(form.FirstName, validators.NameRX), "firstName", "Nome inválido")
	form.CheckField(validators.Matches(form.LastName, validators.NameRX), "lastName", "Nome inválido")
	form.CheckField(validators.Matches(form.Email, validators.EmailRX), "email", "Email inválido")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusSeeOther, "profile.tmpl.html", data)
		return
	}

	user := models.User{
		FullName:  form.FullName,
		FirstName: form.FirstName,
		LastName:  form.LastName,
		Email:     form.Email,
		Phone:     form.Phone,
	}

	id := app.sessionManager.GetInt(r.Context(), userIDKey)
	app.logger.Info("Atualizando perfil de usuário...", userIDKey, id)
	err := app.users.Update(id, user)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "successFlash", "Seu perfil foi atualizado com sucesso!")
	http.Redirect(w, r, "/meu-perfil", http.StatusSeeOther)
}

func (app *application) cartView(w http.ResponseWriter, r *http.Request) {
	productColors, err := app.productColors.GetAll()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	productSizes, err := app.productSizes.GetAll()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := app.newTemplateData(r)
	data.ProductColors = productColors
	data.ProductSizes = productSizes
	app.render(w, r, http.StatusOK, "cart.tmpl.html", data)
}

func (app *application) changePasswordView(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = models.ChangePasswordForm{}
	app.render(w, r, http.StatusOK, "change-password.tmpl.html", data)
}

func (app *application) changePasswordPost(w http.ResponseWriter, r *http.Request) {
	var form models.ChangePasswordForm

	if err := app.decodePostForm(r, &form); err != nil {
		app.serverError(w, r, err)
		return
	}

	form.CheckField(validators.NotBlank(form.CurrentPassword), "currentPassword", "Preencha este campo")
	form.CheckField(validators.NotBlank(form.NewPassword), "newPassword", "Preencha este campo")
	form.CheckField(validators.NotBlank(form.ConfirmPassword), "confirmPassword", "Preencha este campo")

	form.CheckField(validators.MinChars(form.NewPassword, 8), "newPassword", "Sua nova senha deve conter pelo menos 8 caracteres")
	form.CheckField(validators.IsEqual(form.NewPassword, form.ConfirmPassword), "confirmPassword", "As senhas devem ser iguais")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusSeeOther, "change-password.tmpl.html", data)
		return
	}

	id := app.sessionManager.GetInt(r.Context(), userIDKey)
	app.logger.Info("Atualizando senha de usuário...", userIDKey, id)
	err := app.users.ChangePassword(id, form.CurrentPassword, form.NewPassword)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Senha atual inválida")
		} else {
			form.AddNonFieldError(err.Error())
		}

		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusSeeOther, "change-password.tmpl.html", data)
		return
	}

	app.sessionManager.Put(r.Context(), "successFlash", "Sua senha foi atualizada com sucesso!")
	http.Redirect(w, r, "/meu-perfil/alterar-senha", http.StatusSeeOther)
}

func (app *application) wishlistView(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.render(w, r, http.StatusOK, "wishlist.tmpl.html", data)
}

func (app *application) singleProductView(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["slug"]

	// Requisitar da view
	product, err := app.products.GetFromView(slug)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			app.notFound(w, r)
			return
		}
		app.serverError(w, r, err)
		return
	}

	// Requisitar da view
	productImages, err := app.productImage.GetAllByProduct(int(product.ID))
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	productSizes, err := app.productSizes.GetAllByProduct(int(product.SizeCategoryID))
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	productColors, err := app.productColors.GetAll()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := app.newTemplateData(r)
	data.Form = models.ShippingForm{}
	data.ProductView = product
	data.ProductImages = productImages
	data.ProductSizes = productSizes
	data.ProductColors = productColors

	app.render(w, r, http.StatusOK, "single-product.tmpl.html", data)
}

func (app *application) singleCategoryView(w http.ResponseWriter, r *http.Request) {
	var vars = mux.Vars(r)
	category := vars["slug"]

	products, err := app.products.GetFromCategoryParent(category)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := app.newTemplateData(r)
	data.ProductsView = products
	app.render(w, r, http.StatusOK, "single-category.tmpl.html", data)
}

func (app *application) deleteAccountView(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = models.DeleteAccountForm{}
	app.render(w, r, http.StatusOK, "delete-account.tmpl.html", data)
}

func (app *application) deleteAccountPost(w http.ResponseWriter, r *http.Request) {
	err := app.users.PermaDelete(app.sessionManager.GetInt(r.Context(), userIDKey))

	if err != nil {
		var form models.DeleteAccountForm
		form.AddNonFieldError("Erro ao excluir conta. Por favor, tente novamente mais tarde ou entre em contato com o suporte")
		data := app.newTemplateData(r)
		data.Form = form

		app.render(w, r, http.StatusOK, "delete-account.tmpl.html", data)
		return
	}

	app.sessionManager.Remove(r.Context(), "userID")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) addressesView(w http.ResponseWriter, r *http.Request) {
	addresses, err := app.addresses.GetAll(app.sessionManager.GetInt(r.Context(), userIDKey))

	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := app.newTemplateData(r)
	data.Addresses = addresses

	app.render(w, r, http.StatusOK, "addresses.tmpl.html", data)
}

func (app *application) newAddressView(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = models.AddressForm{}
	app.render(w, r, http.StatusOK, "new-address.tmpl.html", data)
}

func (app *application) newAddressPost(w http.ResponseWriter, r *http.Request) {
	var form models.AddressForm

	if err := app.decodePostForm(r, &form); err != nil {
		app.serverError(w, r, err)
		return
	}

	form.CheckField(validators.NotBlank(form.FullName), "fullName", "Preencha o nome corretamente")
	form.CheckField(validators.NotBlank(form.CEP), "CEP", "Preencha o CEP corretamente")
	form.CheckField(validators.NotBlank(form.Street), "street", "Preencha o nome da rua corretamente")
	form.CheckField(validators.NotBlank(form.Number), "number", "Preencha o número da residência corretamente")
	form.CheckField(validators.NotBlank(form.District), "district", "Preencha o bairro corretamente")
	form.CheckField(validators.NotBlank(form.City), "city", "Preencha a cidade corretamente")
	form.CheckField(validators.NotBlank(form.UF), "uf", "Preencha a UF corretamente")

	form.CheckField(validators.Matches(form.FullName, validators.NameRX), "fullName", "Nome inválido")
	form.CheckField(validators.Matches(form.CEP, validators.NumberRX), "cep", "CEP inválido")
	form.CheckField(validators.Matches(form.Number, validators.NumberRX), "number", "Número inválido")
	form.CheckField(validators.Matches(form.District, validators.NameRX), "district", "Bairro inválido")
	form.CheckField(validators.Matches(form.City, validators.NameRX), "city", "Cidade inválida")
	form.CheckField(validators.Matches(form.UF, validators.NameRX), "uf", "UF inválida")

	if form.Complement != "" {
		form.CheckField(validators.Matches(form.Complement, validators.NameRX), "complement", "Nome inválido")
	}

	if form.Reference != "" {
		form.CheckField(validators.Matches(form.Reference, validators.NameRX), "reference", "Nome inválido")
	}

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusSeeOther, "new-address.tmpl.html", data)
		return
	}

	// Criar endereço
	address := models.Address{
		Name:       form.FullName,
		Cep:        form.CEP,
		Street:     form.Street,
		Number:     form.Number,
		Complement: form.Complement,
		District:   form.District,
		City:       form.City,
		UF:         form.UF,
		Reference:  form.Reference,
	}

	app.logger.Info("Criando endereço...")
	addressID, err := app.addresses.Create(address)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// Criar hashID de 101 caracteres
	hashID, err := app.generateHashID(addressID, 101)

	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// Adicionar hashID
	app.addresses.AddHashID(addressID, hashID)

	// Criar associação entre o endereço e o usuário
	// Pegar o ID do usuário na sessão
	userID := app.sessionManager.GetInt(r.Context(), userIDKey)
	err = app.userAddress.Create(userID, addressID)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.logger.Info("Endereço criado com sucesso!", "userID", userID, "addressID", addressID)
	http.Redirect(w, r, "/meu-perfil/enderecos", http.StatusSeeOther)
}

func (app *application) editAddressView(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	hashID := params["hashID"]

	address, err := app.addresses.GetByHashID(hashID)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := app.newTemplateData(r)
	data.Address = address
	data.Form = models.AddressForm{
		FullName:   address.Name,
		CEP:        address.Cep,
		Street:     address.Street,
		Number:     address.Number,
		Complement: address.Complement,
		District:   address.District,
		City:       address.City,
		UF:         address.UF,
		Reference:  address.Reference,
	}

	app.render(w, r, http.StatusOK, "edit-address.tmpl.html", data)
}

func (app *application) editAddressPost(w http.ResponseWriter, r *http.Request) {
	var form models.AddressForm

	if err := app.decodePostForm(r, &form); err != nil {
		app.serverError(w, r, err)
		return
	}

	form.CheckField(validators.NotBlank(form.FullName), "fullName", "Preencha o nome corretamente")
	form.CheckField(validators.NotBlank(form.CEP), "CEP", "Preencha o CEP corretamente")
	form.CheckField(validators.NotBlank(form.Street), "street", "Preencha o nome da rua corretamente")
	form.CheckField(validators.NotBlank(form.Number), "number", "Preencha o número da residência corretamente")
	form.CheckField(validators.NotBlank(form.District), "district", "Preencha o bairro corretamente")
	form.CheckField(validators.NotBlank(form.City), "city", "Preencha a cidade corretamente")
	form.CheckField(validators.NotBlank(form.UF), "uf", "Preencha o estado corretamente")

	form.CheckField(validators.Matches(form.FullName, validators.NameRX), "fullName", "Nome inválido")
	form.CheckField(validators.Matches(form.CEP, validators.NumberRX), "cep", "CEP inválido")
	form.CheckField(validators.Matches(form.Number, validators.NumberRX), "number", "Número inválido")
	form.CheckField(validators.Matches(form.District, validators.NameRX), "district", "Bairro inválido")
	form.CheckField(validators.Matches(form.City, validators.NameRX), "city", "Cidade inválida")
	form.CheckField(validators.Matches(form.UF, validators.NameRX), "state", "UF inválida")

	if form.Complement != "" {
		form.CheckField(validators.Matches(form.Complement, validators.NameRX), "complement", "Nome inválido")
	}

	if form.Reference != "" {
		form.CheckField(validators.Matches(form.Reference, validators.NameRX), "reference", "Nome inválido")
	}

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusSeeOther, "new-address.tmpl.html", data)
		return
	}

	// Atualizar endereço
	address := models.Address{
		Name:       form.FullName,
		Cep:        form.CEP,
		Street:     form.Street,
		Number:     form.Number,
		Complement: form.Complement,
		District:   form.District,
		City:       form.City,
		UF:         form.UF,
		Reference:  form.Reference,
	}

	params := mux.Vars(r)
	hashID := params["hashID"]

	app.logger.Info("Atualizando endereço...")
	err := app.addresses.Update(hashID, address)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "successFlash", "Endereço atualizado com sucesso!")
	http.Redirect(w, r, "/meu-perfil/enderecos", http.StatusSeeOther)
}

func (app *application) deleteAddressPost(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	hashID := params["hashID"]

	err := app.addresses.PermaDeleteByHashID(hashID)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "successFlash", "Endereço excluído com sucesso!")
	http.Redirect(w, r, "/meu-perfil/enderecos", http.StatusSeeOther)
}
