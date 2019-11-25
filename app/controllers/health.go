package controllers

import (
	"net/http"
	"runtime"

	"github.com/briansimoni/stereodose/app/models"

	"github.com/briansimoni/stereodose/app/util"
)

// Health is the struct that ultimately turns into JSON output from this controller
type Health struct {
	GoRoutines          int
	CPUs                int
	HealthyDBConnection bool
}

// HealthController has methods to self-examine the server's health
type HealthController struct {
	db *models.StereoDoseDB
}

// NewHealthController returns a pointer to a HealthController
func NewHealthController(db *models.StereoDoseDB) *HealthController {
	controller := &HealthController{
		db: db,
	}
	controller.db = db
	return controller
}

// CheckHealth is an AppHandler that returns the status of the app's health
func (h *HealthController) CheckHealth(w http.ResponseWriter, r *http.Request) error {
	healthyConnection := true
	err := h.db.DB.DB().Ping()
	if err != nil {
		healthyConnection = false
	}
	health := Health{
		GoRoutines:          runtime.NumGoroutine(),
		CPUs:                runtime.NumCPU(),
		HealthyDBConnection: healthyConnection,
	}
	if !healthyConnection {
		w.WriteHeader(http.StatusInternalServerError)
	}
	util.JSON(w, health)
	return nil
}
