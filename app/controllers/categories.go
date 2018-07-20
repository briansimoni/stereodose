package controllers

import (
	"net/http"

	"github.com/briansimoni/stereodose/app/models"
	"github.com/briansimoni/stereodose/app/util"
	"github.com/pkg/errors"
)

// CategoriesController contains all of the Handler functions related to
// playlist categories. Like Weed/Chill
type CategoriesController struct{}

// GetAvailableCategories sends a JSON object with available categories to use
// The server will deny requests to tag playlists with invalid categories
func (c *CategoriesController) GetAvailableCategories(w http.ResponseWriter, r *http.Request) error {
	err := util.JSON(w, models.Categories)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
