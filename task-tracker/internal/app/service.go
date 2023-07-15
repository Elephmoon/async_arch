package app

import (
	"net/http"
	"task-tracker/internal/pkg/usecases/main_page"

	"github.com/go-chi/chi/v5"
)

type MainPageHandler interface {
	HandleMainRequest(resp http.ResponseWriter, _ *http.Request)
}

type Service struct {
	Router          chi.Router
	mainPageHandler MainPageHandler
}

func SetupHTTPService() Service {
	service := Service{
		Router:          chi.NewRouter(),
		mainPageHandler: main_page.New(),
	}
	service.setupRouters()

	return service
}

func (h *Service) setupRouters() {
	h.Router.Route("/", func(r chi.Router) {
		r.Get("/", h.mainPageHandler.HandleMainRequest)
	})
}
