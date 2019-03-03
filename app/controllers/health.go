package controllers

import (
	"net/http"
	"runtime"

	"github.com/briansimoni/stereodose/app/util"
)

type Health struct {
	GoRoutines int
	CPUs       int
}

type HealthController struct{}

func NewHealthController() *HealthController {
	return new(HealthController)
}

func (h *HealthController) CheckHealth(w http.ResponseWriter, r *http.Request) error {
	health := Health{
		GoRoutines: runtime.NumGoroutine(),
		CPUs:       runtime.NumCPU(),
	}
	util.JSON(w, health)
	return nil
}
