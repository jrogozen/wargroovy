package utils

import (
	"bytes"
	"fmt"
	"github.com/go-chi/render"
	"net/http"
)

func Message(status bool, message string) map[string]interface{} {
	return map[string]interface{}{"status": status, "message": message}
}

func Respond(w http.ResponseWriter, r *http.Request, data map[string]interface{}) {
	render.JSON(w, r, data)
}

func StringWithDefault(val, fallback string) string {
	if val != "" {
		return val
	} else {
		return fallback
	}
}

func MapToString(m map[string]interface{}) string {
	b := new(bytes.Buffer)
	for key, value := range m {
		fmt.Fprintf(b, "%s=\"%s\"\n", key, value)
	}
	return b.String()
}
