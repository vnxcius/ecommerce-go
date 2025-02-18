package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gorilla/mux"
	"github.com/vnxcius/ecommerce-go/internal/models"
	"github.com/vnxcius/ecommerce-go/internal/validators"
	"gorm.io/gorm"
)

/*
Os handlers das páginas são divididos em View, Create,
Post, Update e Delete como sufixo de cada função.

View: são handlers para exibir informações ao usuário
nas telas.

Post: são handlers para processar dados enviados pelo
usuário.

Update: são handlers para processar dados enviados pelo
usuário, geralmente para atualizar registros.

Delete: são handlers para deletar dados.
*/

func (app *application) createProductPost(w http.ResponseWriter, r *http.Request) {
	var form models.ProductForm

	err := app.decodeMultipartPostForm(r, &form)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	form.CheckField(validators.NotBlank(form.Name), "name", "Preencha o nome do produto")
	form.CheckField(validators.NotBlank(form.Description), "description", "Preencha a descricão do produto")
	form.CheckField(validators.NotBlankInt(form.Category), "category", "Selecione uma categoria")
	form.CheckField(validators.NotBlankInt(form.Size), "size", "Selecione um tamanho")
	form.CheckField(validators.NotBlankInt(form.Color), "color", "Selecione uma cor")
	form.CheckField(validators.NotBlankInt(form.Price), "price", "Defina o preço do produto")
	form.CheckField(validators.MaxChars(form.Name, 255), "name", "O nome deve ter menos de 255 caracteres")

	if !form.Valid() {
		categories, _ := app.productCategories.GetAll()
		colors, _ := app.productColors.GetAll()
		sizes, _ := app.productSizes.GetAll()
		sizeCategories, _ := app.productSizeCategories.GetAll()

		data := app.newTemplateData(r)
		data.Form = form
		data.ProductCategories = categories
		data.ProductColors = colors
		data.ProductSizes = sizes
		data.ProductSizeCategories = sizeCategories
		app.render(w, r, http.StatusUnprocessableEntity, "product-create.tmpl.html", data)
		return
	}

	// Criar produto
	app.logger.Info("Criando novo produto...", "name", form.Name)
	product, err := app.products.Create(form.Name, form.Description, form.Category)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// Criar hashID para código SKU do produto
	hashID, _ := app.generateHashID(product.ID, 8)

	// Criar item de produto
	productItem := models.ProductItem{
		ProductID: product.ID,
		Price:     form.Price,
		Discount:  form.Discount,
		Code:      hashID,
		Active:    form.Active,
		ColorID:   form.Color,
	}

	app.logger.Info("Criando novo item de produto...", "productID", product.ID)
	productItemID, err := app.productItem.Create(productItem)

	if err != nil {
		app.logger.Error("Falha ao criar item de produto. Deletando produto...", "productID", product.ID)
		subErr := app.products.PermaDelete(product.ID)
		if subErr != nil {
			app.logger.Error("Falha ao deletar o produto.", "error", subErr)
		}
		app.serverError(w, r, err)
		return
	}

	// Criar variáveis do produto (tamanho e estoque)
	app.logger.Info("Criando novas variáveis produto...", "productID", product.ID)
	err = app.productVariation.Create(productItemID, uint(form.Size), form.InStock)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// Pegar arquivo de imagem do produto
	// file, _, err := r.FormFile("image")

	// if err != nil {
	// 	app.serverError(w, r, err)
	// 	return
	// }
	// defer file.Close()

	// Upar para o cloudinary
	// imageName, err := app.uploadImage(file, "product")
	// if err != nil {
	// 	app.logger.Info("Erro ao salvar imagem no cloudinary, deletando imagem do banco de dados...", "productID", product.ID)
	// 	app.products.PermaDelete(product.ID)

	// 	app.serverError(w, r, err)
	// 	return
	// }

	app.logger.Info("Salvando imagem no banco de dados...", "productID", product.ID)
	// err = app.productImage.Create(imageName, productItemID)
	err = app.productImage.Create("product-placeholder.webp", productItemID)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// err = app.uploadMultipleImages(r.MultipartForm.File["images"], product.ID, productItemID)
	// if err != nil {
	// 	app.serverError(w, r, err)
	// 	return
	// }

	app.sessionManager.Put(r.Context(), "successFlash", "Produto criado com sucesso!")
	http.Redirect(w, r, "/produtos/cadastrar", http.StatusSeeOther)
}

func (app *application) categoriesPost(w http.ResponseWriter, r *http.Request) {
	var form models.CategoryForm

	err := app.decodeMultipartPostForm(r, &form)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// Validação
	form.CheckField(validators.NotBlank(form.Name), "name", "Preencha o nome da categoria")
	form.CheckField(validators.MaxChars(form.Name, 50), "name", "O nome deve ter menos de 50 caracteres")
	form.CheckField(validators.MaxChars(form.Description, 255), "description", "A descricão deve ter menos de 255 caracteres")

	if !form.Valid() {
		data := app.newTemplateData(r)
		categories, _ := app.productCategories.GetAll()
		data.ProductCategories = categories
		data.Form = form

		app.render(w, r, http.StatusUnprocessableEntity, "product-categories.tmpl.html", data)
		return
	}

	// Converter para int, caso exista
	parentCategory := 0
	if form.Parent != "" {
		var err error
		parentCategory, err = strconv.Atoi(form.Parent)
		if err != nil {
			app.serverError(w, r, err)
			return
		}
	}

	// Pegar arquivo
	file, _, err := r.FormFile("image")

	if err != nil && err != http.ErrMissingFile {
		app.logger.Error("ERRO 3", "error", err)
		app.serverError(w, r, err)
		return
	}

	if file != nil {
		defer file.Close()
		// Upa a imagem para o cloudinary.
		// Recebe arquivo e qual o contexto do upload, neste caso
		// é uma imagem de categoria.
		imageName, err := app.uploadImage(file, "category")

		if err != nil {
			app.serverError(w, r, err)
			return
		}

		app.logger.Info("Criando categoria...", "category", form.Name)
		err = app.productCategories.Create(form.Name, form.Description, imageName, parentCategory)
		if err != nil {
			app.serverError(w, r, err)
			return
		}

		app.logger.Info("Categoria "+form.Name+" criada com sucesso!", "category", form.Name)
		app.sessionManager.Put(r.Context(), "successFlash", "Categora criada com sucesso!")
		http.Redirect(w, r, "/produtos/categorias", http.StatusSeeOther)
		return
	}

	app.logger.Info("Criando categoria...", "category", form.Name)
	err = app.productCategories.Create(form.Name, form.Description, "", parentCategory)

	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.logger.Info("Categoria "+form.Name+" criada com sucesso!", "category", form.Name)
	app.sessionManager.Put(r.Context(), "successFlash", "Categora criada com sucesso!")
	http.Redirect(w, r, "/produtos/categorias", http.StatusSeeOther)
}

func (app *application) categoryUpdate(w http.ResponseWriter, r *http.Request) {
	// TODO
}

func (app *application) categoryDelete(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	categoryID := params["category_id"]
	image := params["image"]

	categoryIDInt, err := strconv.Atoi(categoryID)

	if err != nil {
		app.serverError(w, r, err)
		return
	}

	if image != "null" {
		app.logger.Info("Deletando imagem da categoria...", "image", image)
		res, _ := app.cld.Upload.Destroy(app.ctx, uploader.DestroyParams{
			PublicID: image,
		})

		if res.Result != "ok" {
			app.sessionManager.Put(r.Context(), "errorFlash", "Erro: "+res.Result)
			http.Redirect(w, r, "/produtos/categorias", http.StatusSeeOther)
			return
		}

		app.logger.Info("Imagem da categoria deletada com sucesso!", "image", image)
	}

	app.logger.Info("Deletando categoria...", "category_id", categoryID)
	err = app.productCategories.PermaDelete(categoryIDInt)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.logger.Info("Categoria deletada com sucesso!", "category_id", categoryID)
	app.sessionManager.Put(r.Context(), "successFlash", "Categoria deletada com sucesso!")
	http.Redirect(w, r, "/produtos/categorias", http.StatusSeeOther)
}

func (app *application) sizePost(w http.ResponseWriter, r *http.Request) {
	var form models.SizeForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	form.CheckField(validators.NotBlank(form.NameSize), "nameSize", "Preencha o nome do tamanho")
	form.CheckField(validators.NotBlank(form.SortOrder), "sortOrder", "Defina a ordem de exibição")
	form.CheckField(validators.NotBlank(form.Parent), "parent", "Selecione a qual categoria este tamanho pertence")

	form.CheckField(validators.MaxChars(form.NameSize, 10), "nameSize", "O nome deve ter menos de 10 caracteres")
	form.CheckField(validators.MaxChars(form.SortOrder, 2), "sortOrder", "Insira apenas números de 1 a 99")

	if !form.Valid() {
		sizes, _ := app.productSizes.GetAll()
		sizeCategories, _ := app.productSizeCategories.GetAll()

		data := app.newTemplateData(r)
		data.ProductSizes = sizes
		data.ProductSizeCategories = sizeCategories
		data.Form = form

		app.render(w, r, http.StatusUnprocessableEntity, "product-sizes.tmpl.html", data)
		return
	}

	// Converter sortOrder para int
	sortOrder, err := strconv.Atoi(form.SortOrder)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// Converter parent para int
	parent, err := strconv.Atoi(form.Parent)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.logger.Info("Criando tamanho...", "size", form.Name)
	err = app.productSizes.Create(form.NameSize, int32(sortOrder), int32(parent))

	if err != nil {
		form.AddNonFieldError(err.Error())
		sizes, _ := app.productSizes.GetAll()
		sizeCategories, _ := app.productSizeCategories.GetAll()

		data := app.newTemplateData(r)
		data.ProductSizes = sizes
		data.ProductSizeCategories = sizeCategories
		data.Form = form

		app.render(w, r, http.StatusUnprocessableEntity, "product-sizes.tmpl.html", data)
		return
	}

	app.sessionManager.Put(r.Context(), "successFlash", "Tamanho criado com sucesso!")
	app.logger.Info("Tamanho "+form.Name+" criado com sucesso!", "size", form.Name)
	http.Redirect(w, r, "/produtos/tamanhos", http.StatusSeeOther)
}

func (app *application) sizeCategoriesPost(w http.ResponseWriter, r *http.Request) {
	var form models.SizeForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	form.CheckField(validators.NotBlank(form.Name), "name", "Preencha o nome da categoria")
	form.CheckField(validators.MaxChars(form.Name, 25), "name", "O nome deve ter menos de 25 caracteres")

	if !form.Valid() {
		sizes, _ := app.productSizes.GetAll()
		sizeCategories, _ := app.productSizeCategories.GetAll()

		data := app.newTemplateData(r)
		data.ProductSizes = sizes
		data.ProductSizeCategories = sizeCategories
		data.Form = form

		app.render(w, r, http.StatusUnprocessableEntity, "product-sizes.tmpl.html", data)
		return
	}

	app.logger.Info("Criando categoria de tamanho...", "size", form.Name)
	err = app.productSizeCategories.Create(form.Name)

	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			form.AddNonFieldError("Categoria de tamanho já existente")
		} else {
			form.AddNonFieldError(err.Error())
		}

		sizes, _ := app.productSizes.GetAll()
		sizeCategories, _ := app.productSizeCategories.GetAll()

		data := app.newTemplateData(r)
		data.ProductSizes = sizes
		data.ProductSizeCategories = sizeCategories
		data.Form = form

		app.render(w, r, http.StatusUnprocessableEntity, "product-sizes.tmpl.html", data)
		return
	}

	app.sessionManager.Put(r.Context(), "successFlash", "Categoria de tamanho criada com sucesso!")
	app.logger.Info("Categoria de tamanho "+form.Name+" criada com sucesso!", "size", form.Name)
	http.Redirect(w, r, "/produtos/tamanhos", http.StatusSeeOther)
}

func (app *application) sizeDelete(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	sizeID := params["size_id"]

	sizeIDInt, err := strconv.Atoi(sizeID)

	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.logger.Info("Deletando tamanho...", "size_id", sizeID)
	err = app.productSizes.SoftDelete(sizeIDInt)

	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.logger.Info("Tamanho deletado com sucesso!", "size_id", sizeID)
	app.sessionManager.Put(r.Context(), "successFlash", "Tamanho deletado com sucesso!")
	http.Redirect(w, r, "/produtos/tamanhos", http.StatusSeeOther)
}

func (app *application) sizeCategoriesDelete(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	sizeCategoryID := params["size_category_id"]

	sizeCategoryIDInt, err := strconv.Atoi(sizeCategoryID)

	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.logger.Info("Deletando categoria de tamanho...", "size_category_id", sizeCategoryID)
	err = app.productSizeCategories.SoftDelete(sizeCategoryIDInt)

	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.logger.Info("Categoria de tamanho deletada com sucesso!", "size_category_id", sizeCategoryID)
	app.sessionManager.Put(r.Context(), "successFlash", "Categoria de tamanho deletada com sucesso!")
	http.Redirect(w, r, "/produtos/tamanhos", http.StatusSeeOther)
}

func (app *application) sizeUpdate(w http.ResponseWriter, r *http.Request) {
	var form models.SizeForm
	params := mux.Vars(r)

	err := app.decodePostForm(r, &form)

	if err != nil {
		app.serverError(w, r, err)
		return
	}

	form.CheckField(validators.NotBlank(form.Name), "sizeName", "Preencha o nome do tamanho")
	form.CheckField(validators.MaxChars(form.Name, 10), "sizeName", "O nome deve ter menos de 10 caracteres")
	form.CheckField(validators.NotBlank(form.SortOrder), "sortOrder", "Defina a ordem de exibição")
	form.CheckField(validators.MaxChars(form.SortOrder, 2), "sortOrder", "Insira apenas números de 1 a 99")

	if !form.Valid() {
		data := app.newTemplateData(r)

		sizes, _ := app.productSizes.GetAll()
		sizeCategories, _ := app.productSizeCategories.GetAll()
		data.ProductSizes = sizes
		data.ProductSizeCategories = sizeCategories
		data.Form = form

		app.render(w, r, http.StatusUnprocessableEntity, "product-sizes.tmpl.html", data)
		return
	}

	// Converter para int
	id, _ := strconv.Atoi(params["size_id"])
	sortOrder, _ := strconv.Atoi(form.SortOrder)

	app.logger.Info("Atualizando tamanho...", "size_id", id)
	err = app.productSizes.Update(id, form.Name, int32(sortOrder))

	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "successFlash", "Tamanho atualizado com sucesso!")
	app.logger.Info("Tamanho "+form.Name+" atualizado com sucesso!", "size_id", id)
	http.Redirect(w, r, "/produtos/tamanhos", http.StatusSeeOther)
}

func (app *application) sizeCategoriesUpdate(w http.ResponseWriter, r *http.Request) {
	var form models.SizeForm
	params := mux.Vars(r)

	err := app.decodePostForm(r, &form)

	if err != nil {
		app.serverError(w, r, err)
		return
	}

	form.CheckField(validators.NotBlank(form.Name), "name", "Preencha o nome da categoria")
	form.CheckField(validators.MaxChars(form.Name, 25), "name", "O nome deve ter menos de 50 caracteres")

	if !form.Valid() {
		data := app.newTemplateData(r)

		sizes, _ := app.productSizes.GetAll()
		sizeCategories, _ := app.productSizeCategories.GetAll()
		data.ProductSizes = sizes
		data.ProductSizeCategories = sizeCategories
		data.Form = form

		app.render(w, r, http.StatusUnprocessableEntity, "product-sizes.tmpl.html", data)
		return
	}

	// Converter ID para int
	id, _ := strconv.Atoi(params["size_category_id"])

	app.logger.Info("Atualizando categoria de tamanho...", "size_category_id", id)
	err = app.productSizeCategories.Update(id, form.Name)

	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "successFlash", "Categoria de tamanho atualizada com sucesso!")
	app.logger.Info("Categoria de tamanho atualizada com sucesso!", "size_category_id", id)
	http.Redirect(w, r, "/produtos/tamanhos", http.StatusSeeOther)
}

func (app *application) colorsPost(w http.ResponseWriter, r *http.Request) {
	var form models.ColorForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	form.CheckField(validators.NotBlank(form.Name), "name", "Preencha o nome da cor")
	form.CheckField(validators.MaxChars(form.Name, 20), "name", "O nome deve ter menos de 20 caracteres")
	form.CheckField(validators.NotBlank(form.Hex), "hex", "Selecione o tom da cor")

	if !form.Valid() {
		data := app.newTemplateData(r)
		colors, _ := app.productColors.GetAll()
		data.ProductColors = colors
		data.Form = form

		app.render(w, r, http.StatusUnprocessableEntity, "product-colors.tmpl.html", data)
		return
	}

	app.logger.Info("Criando nova cor...", "name", form.Name)
	err = app.productColors.Create(form.Name, form.Hex)

	if err != nil {
		app.sessionManager.Put(r.Context(), "errorFlash", "Erro ao criar cor: "+err.Error())

		if err == gorm.ErrDuplicatedKey {
			app.sessionManager.Put(r.Context(), "errorFlash", "Erro ao criar cor: esta cor já existe (nome ou código da cor)")
		}
		http.Redirect(w, r, "/produtos/cores", http.StatusSeeOther)
		return
	}

	app.sessionManager.Put(r.Context(), "successFlash", "Cor criada com sucesso!")
	http.Redirect(w, r, "/produtos/cores", http.StatusSeeOther)
}

func (app *application) colorsUpdate(w http.ResponseWriter, r *http.Request) {
	var form models.ColorForm
	params := mux.Vars(r)

	err := app.decodePostForm(r, &form)

	if err != nil {
		app.serverError(w, r, err)
		return
	}

	form.CheckField(validators.NotBlank(form.Name), "name", "Preencha o nome da cor")
	form.CheckField(validators.MaxChars(form.Name, 15), "name", "O nome deve ter menos de 25 caracteres")

	if !form.Valid() {
		data := app.newTemplateData(r)

		colors, _ := app.productColors.GetAll()
		data.ProductColors = colors
		data.Form = form

		app.render(w, r, http.StatusUnprocessableEntity, "product-colors.tmpl.html", data)
		return
	}

	// Converter ID para int
	id, _ := strconv.Atoi(params["color_id"])

	app.logger.Info("Atualizando cor...", "color_id", id)
	err = app.productColors.Update(id, form.Name)

	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "successFlash", "Cor atualizada com sucesso!")
	app.logger.Info("Cor atualizada com sucesso!", "color_id", id)
	http.Redirect(w, r, "/produtos/cores", http.StatusSeeOther)
}

func (app *application) colorsDelete(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id, _ := strconv.Atoi(params["color_id"])

	app.logger.Info("Excluindo cor...", "color_id", id)
	err := app.productColors.Delete(id)

	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.logger.Info("Cor excluída com sucesso!", "color_id", id)
	app.sessionManager.Put(r.Context(), "successFlash", "Cor excluída com sucesso!")
	http.Redirect(w, r, "/produtos/cores", http.StatusSeeOther)
}
