package API

import (
	"github.com/ghuvrons/fota-server-go/API/controllers"

	"github.com/go-chi/chi/v5"
)

func route(r *chi.Mux) {
	r.Post("/upload-firmware", controllers.UploadFirmware)
}
