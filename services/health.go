package services

import (
	"github.com/jinzhu/gorm"
	"github.com/sample-full-api/utils"
	"github.com/sample-full-api/views"
	"github.com/sirupsen/logrus"
)

type HealthService interface {
	HealthCheck() (views.HealthCheckResponse, []error)
}

type healthService struct {
	name   string
	db     *gorm.DB
	logger *logrus.Logger
}

func NewHealthService(deps utils.Dependencies) *healthService {
	return &healthService{
		name:   "forecast service",
		db:     deps.Db,
		logger: deps.Logger,
	}
}

func (h *healthService) HealthCheck() (response views.HealthCheckResponse, errs []error) {
	response.Name = h.name
	response.Database = true

	if err := h.dbHealthCheck(); err != nil {
		errs = append(errs, err)
		h.logger.WithError(err)
		response.Database = false
	}

	return
}

func (h *healthService) dbHealthCheck() error {
	return h.db.DB().Ping()
}
