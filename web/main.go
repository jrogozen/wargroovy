package main

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/jrogozen/wargroovy/handlers/auth"
	"github.com/jrogozen/wargroovy/handlers/maps"
	"github.com/jrogozen/wargroovy/handlers/user"
	"github.com/jrogozen/wargroovy/internal/config"
	log "github.com/sirupsen/logrus"
	"google.golang.org/appengine"
	"net/http"
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
		r.Mount("/map", maps.Routes(configuration))
		r.Mount("/auth", auth.Routes(configuration))
	})

	router.Route("/", func(r chi.Router) {
		r.Get("/", helloHandler())
	})

	return router
}

func helloHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		render.JSON(w, r, map[string]interface{}{"status": "ok"})
	})
}

func main() {
	configuration, err := config.New()

	// does this set for everything?
	log.SetFormatter(&log.TextFormatter{
		ForceColors: true,
	})

	if err != nil {
		log.Warn("Error reading configuration", err)
	}

	router := Routes(configuration)

	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		log.Printf("[%s]: %s", method, route)

		return nil
	}

	if err := chi.Walk(router, walkFunc); err != nil {
		log.Panicf("Logging err: %s\n", err.Error())
	}

	log.WithFields(log.Fields{
		"PORT": configuration.Constants.PORT,
	}).Info("Serving application")

	log.Fatal(http.ListenAndServe(":"+configuration.Constants.PORT, router))

	appengine.Main()
}
