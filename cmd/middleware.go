package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/justinas/nosurf"
)

func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy",
			"default-src 'self' blob: cdn.jsdelivr.net res.cloudinary.com; "+
				"style-src 'self' 'unsafe-inline' fonts.googleapis.com cdn.jsdelivr.net; "+
				"script-src 'self' 'unsafe-inline' cdn.jsdelivr.net; "+
				"font-src 'self' fonts.gstatic.com;"+
				"connect-src 'self' viacep.com.br")

		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")
		w.Header().Set("Cache-Control", "public, max-age=2592000")

		next.ServeHTTP(w, r)
	})
}

func (app *application) isMaintenanceMode(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if app.MaintenanceMode {
			w.Header().Set("Cache-Control", "no-store")
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("Estamos em manutenção. Por favor, tente novamente mais tarde."))
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verifica se a URL requisitada começa com /static/
		if strings.HasPrefix(r.URL.Path, "/static/") {
			// Se for uma requisição de arquivo estático, apenas chame o próximo handler
			next.ServeHTTP(w, r)
			return
		}
		
		var ip string
		// Verifica o cabeçalho X-Forwarded-For
		xff := r.Header.Get("X-Forwarded-For")
		if xff != "" {
			// X-Forwarded-For pode conter uma lista de IPs. Pegue o primeiro.
			ips := strings.Split(xff, ",")
			ip = strings.TrimSpace(ips[0])
		} else {
			// Se não houver X-Forwarded-For, verifica o X-Real-IP
			xri := r.Header.Get("X-Real-IP")
			if xri != "" {
				ip = xri
			} else {
				// Caso contrário, use o endereço remoto da requisição
				ip = r.RemoteAddr
			}
		}
		var (
			proto  = r.Proto
			method = r.Method
			uri    = r.URL.RequestURI()
		)

		app.logger.Info("received request", "ip", ip, "proto", proto, "method", method, "uri", uri)

		next.ServeHTTP(w, r)
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a deferred function (which will always be run in the event
		// of a panic as Go unwinds the stack).
		defer func() {
			// Use the builtin recover function to check if there has been a
			// panic or not. If there has...
			if err := recover(); err != nil {
				// Set a "Connection: close" header on the response.
				w.Header().Set("Connection", "close")
				// Call the app.serverError helper method to return a 500
				// Internal Server response.
				app.serverError(w, r, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (app *application) requireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !app.isAuthenticated(r) {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		id := app.sessionManager.GetInt(r.Context(), userIDKey)
		isAdmin, err := app.users.IsAdmin(id)
		if err != nil {
			app.serverError(w, r, err)
			return
		}

		if !isAdmin {
			app.logger.Info("User tried to access admin-only page", "id", id)
			app.logout(w, r)
			return
		}

		// Impede que páginas que requerem que o usuário esteja autenticado
		// sejam armazenadas em cache
		w.Header().Add("Cache-Control", "no-store")
		next.ServeHTTP(w, r)
	})
}

func (app *application) requireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !app.isAuthenticated(r) {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Impede que páginas que requerem que o usuário esteja autenticado
		// sejam armazenadas em cache
		w.Header().Add("Cache-Control", "no-store")
		next.ServeHTTP(w, r)
	})
}

func (app *application) noSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
	})

	return csrfHandler
}

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := app.sessionManager.GetInt(r.Context(), userIDKey)
		if id == 0 {
			app.logger.Info("failed to authenticate user")
			next.ServeHTTP(w, r)
			return
		}

		// otherwise
		exists, err := app.users.Exists(id)
		if err != nil {
			app.logger.Error("failed to check if user exists", "error", err)
			app.serverError(w, r, err)
			return
		}

		if exists {
			ctx := context.WithValue(r.Context(), isAuthenticatedUserKey, true)
			r = r.WithContext(ctx)
		}

		app.logger.Info("user authenticated", "id", id)
		next.ServeHTTP(w, r)
	})
}
