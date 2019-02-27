package photos

import (
	"context"
	"fmt"
	uuid "github.com/gofrs/uuid"
	"github.com/jrogozen/wargroovy/internal/config"
	u "github.com/jrogozen/wargroovy/utils"
	log "github.com/sirupsen/logrus"
	"io"
	"mime/multipart"
	"net/http"
	"path"
)

func Upload(configuration *config.Config, f multipart.File, fh *multipart.FileHeader) (map[string]interface{}, int) {
	// random filename, retaining existing extension.
	name := "map_photos/" + uuid.Must(uuid.NewV4()).String() + path.Ext(fh.Filename)

	log.Info(name)

	ctx := context.Background()
	w := configuration.StorageBucket.Object(name).NewWriter(ctx)

	// Warning: storage.AllUsers gives public read access to anyone.
	// Object policy currently disabled for this bucket
	// w.ACL = []storage.ACLRule{{Entity: storage.AllUsers, Role: storage.RoleReader}}
	w.ContentType = fh.Header.Get("Content-Type")
	w.CacheControl = "public, max-age=86400"

	if _, err := io.Copy(w, f); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Error copying photo")

		return u.Message(false, "Error uploading photo"), http.StatusBadRequest
	}
	if err := w.Close(); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Error closing write stream")

		return u.Message(false, "Error uploading photo"), http.StatusBadRequest
	}

	const publicURL = "https://storage.googleapis.com/%s/%s"
	uploadedURL := fmt.Sprintf(publicURL, configuration.StorageBucketName, name)

	response := u.Message(true, "Uploaded photos")
	response["url"] = uploadedURL

	return response, http.StatusOK
}
