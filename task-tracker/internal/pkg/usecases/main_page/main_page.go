package main_page

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"task-tracker/internal/pkg/assets"
	"task-tracker/internal/pkg/helpers"
	"task-tracker/internal/pkg/models"
)

type MainPage struct {
}

func New() *MainPage {
	return &MainPage{}
}

func (h *MainPage) HandleMainRequest(resp http.ResponseWriter, _ *http.Request) {
	mainPageTemplate, err := template.ParseFS(assets.GetAssets(), "sources/main.html")
	if err != nil {
		errPayload, modelErr := models.NewHttpError(http.StatusInternalServerError, fmt.Errorf("cant load template %w", err))
		if modelErr != nil {
			log.Printf("cant build error %v", err)
		}
		helpers.WriteHttpResponse(resp, http.StatusInternalServerError, errPayload)
		return
	}

	err = mainPageTemplate.Execute(resp, "")
	if err != nil {
		errPayload, modelErr := models.NewHttpError(http.StatusInternalServerError, fmt.Errorf("cant execute template %w", err))
		if modelErr != nil {
			log.Printf("cant build error %v", err)
		}
		helpers.WriteHttpResponse(resp, http.StatusInternalServerError, errPayload)
		return
	}

	resp.WriteHeader(http.StatusOK)
}
