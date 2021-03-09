package controllers

import (
	"picapp/models"
	"picapp/views"
)

func NewGalleries(gs models.GalleryService) *Galleries {
	return &Galleries{
		New:   views.NewView("bootstrap", "galleries/new"),
		gs:        gs,
	}
}

type Galleries struct {
	New   *views.View
	gs        models.GalleryService
}