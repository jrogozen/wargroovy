package photos

import (
	"errors"
	// "github.com/go-chi/render"
	"github.com/jrogozen/wargroovy/internal/config"
	u "github.com/jrogozen/wargroovy/utils"
	"net/http"
)

func UploadPhotos(configuration *config.Config) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if configuration.StorageBucket == nil {
			u.Respond(w, r, u.Message(false, errors.New("storage bucket not defined").Error()))
			return
		}

		r.ParseMultipartForm(32 << 20) // 32MB is the default used by FormFile
		fhs := r.MultipartForm.File["photos"]

		photoURLs := make([]string, 0)

		for _, fh := range fhs {
			f, err := fh.Open()

			resp, status := Upload(configuration, f, fh)

			if resp["status"].(bool) {
				photoURLs = append(photoURLs, resp["url"].(string))
			}

			if err == http.ErrMissingFile {
				w.WriteHeader(status)
				u.Respond(w, r, u.Message(false, "no photo to upload"))
				return
			}

			if err != nil {
				w.WriteHeader(status)
				u.Respond(w, r, u.Message(false, err.Error()))
				return
			}
		}

		response := u.Message(true, "photos uploaded")
		response["urls"] = photoURLs

		w.WriteHeader(http.StatusOK)
		u.Respond(w, r, response)
		return
	})
}
