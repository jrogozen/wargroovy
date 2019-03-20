package photos

import (
	"errors"
	// "github.com/go-chi/render"
	"github.com/jrogozen/wargroovy/internal/config"
	u "github.com/jrogozen/wargroovy/utils"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func UploadPhotos(configuration *config.Config) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if configuration.StorageBucket == nil {
			u.Respond(w, r, u.Message(false, errors.New("storage bucket not defined").Error()))
			return
		}

		err := r.ParseMultipartForm(32 << 20) // 32MB is the default used by FormFile

		if err != nil {
			log.WithField("err", err).Error("error parsing multipart form")
		}

		fhs := r.MultipartForm.File["photos"]

		photoURLs := make([]string, 0)

		log.WithField("fhs", fhs).Info("upload: fhs")

		for _, fh := range fhs {
			f, err := fh.Open()

			resp, status := Upload(configuration, f, fh)

			if resp["status"].(bool) {
				photoURLs = append(photoURLs, resp["url"].(string))
			}

			if err == http.ErrMissingFile {
				log.WithField("err", err).Error("photo upload: missing file in http request")

				w.WriteHeader(status)
				u.Respond(w, r, u.Message(false, "no photo to upload"))
				return
			}

			if err != nil {
				log.WithField("err", err).Error("photo upload: error")

				w.WriteHeader(status)
				u.Respond(w, r, u.Message(false, err.Error()))
				return
			}
		}

		response := u.Message(true, "photos uploaded")
		response["urls"] = photoURLs

		w.WriteHeader(http.StatusOK)

		log.WithFields(log.Fields{
			"response": response,
		}).Info("upload: sending response")

		u.Respond(w, r, response)
		return
	})
}
