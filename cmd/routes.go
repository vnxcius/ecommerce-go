package main

import (
	"io/fs"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	uinterface "github.com/vnxcius/ecommerce-go/interface"
)

func (app *application) routes() http.Handler {
	r := mux.NewRouter().StrictSlash(true) // Aceitar barra no final das URLs
	s := r.Host("{subdomain:admin+}." + app.domain).Subrouter()

	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w, r)
	})

	// Servir arquivos estáticos sem que sejam navegáveis
	staticFS, err := fs.Sub(uinterface.Files, "static")
	if err != nil {
		panic(err)
	}
	fileServer := http.FileServer(http.FS(staticFS))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fileServer))

	// Carregar e salvar dados da sessão
	dynamic := alice.New(app.sessionManager.LoadAndSave, app.noSurf, app.authenticate, app.isMaintenanceMode)
	authenticated := dynamic.Append(app.requireAuth)
	adminProtected := dynamic.Append(app.requireAdmin)

	r.Handle("/ping", dynamic.ThenFunc(app.ping)).Methods("GET")

	// Rotas cliente
	r.Handle("/", dynamic.ThenFunc(app.homeView)).Methods("GET")
	r.Handle("/search", dynamic.ThenFunc(app.searchView)).Methods("GET").Queries("q", "{q}")
	r.Handle("/login", dynamic.ThenFunc(app.loginView)).Methods("GET")
	r.Handle("/login", dynamic.ThenFunc(app.loginPost)).Methods("POST")
	r.Handle("/logout", dynamic.ThenFunc(app.logout)).Methods("POST")
	r.Handle("/registrar", dynamic.ThenFunc(app.registerView)).Methods("GET")
	r.Handle("/registrar", dynamic.ThenFunc(app.registerPost)).Methods("POST")
	r.Handle("/meu-perfil", authenticated.ThenFunc(app.profileView)).Methods("GET")
	r.Handle("/meu-perfil", authenticated.ThenFunc(app.profilePost)).Methods("POST")
	r.Handle("/meu-perfil/lista-de-desejos", authenticated.ThenFunc(app.wishlistView)).Methods("GET")
	r.Handle("/meu-perfil/alterar-senha", authenticated.ThenFunc(app.changePasswordView)).Methods("GET")
	r.Handle("/meu-perfil/alterar-senha", authenticated.ThenFunc(app.changePasswordPost)).Methods("POST")
	r.Handle("/meu-perfil/deletar-conta", authenticated.ThenFunc(app.deleteAccountView)).Methods("GET")
	r.Handle("/meu-perfil/deletar-conta", authenticated.ThenFunc(app.deleteAccountPost)).Methods("POST")

	r.Handle("/meu-perfil/enderecos", authenticated.ThenFunc(app.addressesView)).Methods("GET")
	r.Handle("/meu-perfil/enderecos/new", authenticated.ThenFunc(app.newAddressView)).Methods("GET")
	r.Handle("/meu-perfil/enderecos/new", authenticated.ThenFunc(app.newAddressPost)).Methods("POST")
	r.Handle("/meu-perfil/enderecos/editar/{hashID}", authenticated.ThenFunc(app.editAddressView)).Methods("GET")
	r.Handle("/meu-perfil/enderecos/editar/{hashID}", authenticated.ThenFunc(app.editAddressPost)).Methods("POST")
	r.Handle("/meu-perfil/enderecos/deletar/{hashID}", authenticated.ThenFunc(app.deleteAddressPost)).Methods("POST")

	r.Handle("/p/{slug}", dynamic.ThenFunc(app.singleProductView)).Methods("GET")
	r.Handle("/c/{slug}", dynamic.ThenFunc(app.singleCategoryView)).Methods("GET")
	r.Handle("/carrinho", dynamic.ThenFunc(app.cartView)).Methods("GET")

	// Rotas admin
	s.Handle("/", dynamic.ThenFunc(app.redirectToLogin)).Methods("GET")
	s.Handle("/login", dynamic.ThenFunc(app.loginAdminView)).Methods("GET")
	s.Handle("/login", dynamic.ThenFunc(app.loginAdminPost)).Methods("POST")
	s.Handle("/dashboard", adminProtected.ThenFunc(app.dashboardView)).Methods("GET")
	s.Handle("/logout", adminProtected.ThenFunc(app.logout)).Methods("POST")

	s.Handle("/usuarios", adminProtected.ThenFunc(app.usersView)).Methods("GET")
	s.Handle("/usuarios/registrar-admin", adminProtected.ThenFunc(app.registerAdminView)).Methods("GET")
	s.Handle("/usuarios/registrar-admin", adminProtected.ThenFunc(app.registerAdminPost)).Methods("POST")
	s.Handle("/usuarios/admin/alterar-foto/{user_id}", adminProtected.ThenFunc(app.updateAdminPhoto)).Methods("POST")

	s.Handle("/pedidos", adminProtected.ThenFunc(app.ordersView)).Methods("GET")

	// Rotas de produtos admin
	s.Handle("/produtos", adminProtected.ThenFunc(app.allProductsView)).Methods("GET")
	s.Handle("/produtos/editar/{product_code}", adminProtected.ThenFunc(app.productView)).Methods("GET")
	s.Handle("/produtos/cadastrar", adminProtected.ThenFunc(app.createProductView)).Methods("GET")
	s.Handle("/produtos/cadastrar", adminProtected.ThenFunc(app.createProductPost)).Methods("POST")

	s.Handle("/produtos/categorias", adminProtected.ThenFunc(app.categoriesView)).Methods("GET")
	s.Handle("/produtos/categorias", adminProtected.ThenFunc(app.categoriesPost)).Methods("POST")
	s.Handle("/produtos/categorias/editar/{category_id}/{image}", adminProtected.ThenFunc(app.categoryUpdate)).Methods("POST")
	s.Handle("/produtos/categorias/deletar/{category_id}/{image}", adminProtected.ThenFunc(app.categoryDelete)).Methods("POST")

	s.Handle("/produtos/tamanhos", adminProtected.ThenFunc(app.sizesView)).Methods("GET")
	s.Handle("/produtos/tamanhos", adminProtected.ThenFunc(app.sizePost)).Methods("POST")
	s.Handle("/produtos/tamanhos/editar/{size_id}", adminProtected.ThenFunc(app.sizeUpdate)).Methods("POST")
	s.Handle("/produtos/tamanhos/deletar/{size_id}", adminProtected.ThenFunc(app.sizeDelete)).Methods("POST")

	s.Handle("/produtos/tamanhos/categorias", adminProtected.ThenFunc(app.sizeCategoriesPost)).Methods("POST")
	s.Handle("/produtos/tamanhos/categorias/editar/{size_category_id}", adminProtected.ThenFunc(app.sizeCategoriesUpdate)).Methods("POST")
	s.Handle("/produtos/tamanhos/categorias/deletar/{size_category_id}", adminProtected.ThenFunc(app.sizeCategoriesDelete)).Methods("POST")

	s.Handle("/produtos/cores", adminProtected.ThenFunc(app.colorsView)).Methods("GET")
	s.Handle("/produtos/cores", adminProtected.ThenFunc(app.colorsPost)).Methods("POST")
	s.Handle("/produtos/cores/editar/{color_id}", adminProtected.ThenFunc(app.colorsUpdate)).Methods("POST")
	s.Handle("/produtos/cores/deletar/{color_id}", adminProtected.ThenFunc(app.colorsDelete)).Methods("POST")

	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	return standard.Then(r)
}
