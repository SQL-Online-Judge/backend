package restapi

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
)

func NewRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(contentTypeJSON)

	r.Group(func(r chi.Router) {
		r.Get("/", sayHello("SQL-Online-Judge"))
		r.Post("/login", login)
	})

	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(jwtauth.Authenticator(tokenAuth))
		r.Use(getRole)

		r.Route("/admin", func(r chi.Router) {
			r.Group(func(r chi.Router) {
				r.Use(checkRole("admin"))
				r.Post("/teacher", createTeacher)
			})

			r.Group(func(r chi.Router) {
				r.Use(checkRole("teacher"))
				r.Post("/students", createStudents)
			})
		})
	})

	return r
}
