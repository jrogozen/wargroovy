package utils

import (
	"net/http"

	"github.com/go-chi/render"
)

func Message(status bool, message string) map[string]interface{} {
	return map[string]interface{}{"status": status, "message": message}
}

func Respond(w http.ResponseWriter, r *http.Request, data map[string]interface{}) {
	render.JSON(w, r, data)
}
