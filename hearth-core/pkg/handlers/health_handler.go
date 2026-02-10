package handlers

import (
	"github.com/Mahaveer86619/Hearth/pkg/services"
	"github.com/Mahaveer86619/Hearth/pkg/views"
	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
	hs *services.HealthService
}

func NewHealthHandler(hs *services.HealthService) *HealthHandler {
	return &HealthHandler{hs: hs}
}

func (hh *HealthHandler) Ping(c *gin.Context) {
	resp := &views.Success{}
	resp.SetMessage(hh.hs.Ping()).Send(c)
}

func (hh *HealthHandler) Check(c *gin.Context) {
	resp := &views.Success{}
	resp.SetStatusCode(200).
		SetMessage("Health Check").
		SetData(gin.H{"status": hh.hs.Check()}).
		Send(c)
}
