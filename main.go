package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"

	"github.com/jrogozen/wargroovy/handlers/auth"
	"github.com/jrogozen/wargroovy/handlers/user"
	"github.com/jrogozen/wargroovy/internal/config"
)

func Routes(configuration *config.Config) *chi.Mux {
	router := chi.NewRouter()
	router.Use(
		render.SetContentType(render.ContentTypeJSON), // Set content-Type headers as application/json
		middleware.Logger,          // Log API request calls
		middleware.DefaultCompress, // Compress results, mostly gzipping assets and json
		middleware.RedirectSlashes, // Redirect slashes to no slash URL versions
		middleware.Recoverer,       // Recover from panics without crashing server
	)

	router.Route("/v1/api", func(r chi.Router) {
		r.Mount("/user", user.Routes(configuration))

		r.Mount("/auth", auth.Routes(configuration))
	})

	return router
}

func main() {
	configuration, err := config.New()

	if err != nil {
		log.Panicln("Error reading configuration", err)
	}

	router := Routes(configuration)

	log.Println("Serving application at PORT: " + configuration.Constants.PORT)
	log.Fatal(http.ListenAndServe(":"+configuration.Constants.PORT, router))
}
