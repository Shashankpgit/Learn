package controller

import (
	"app/src/constants"
	"app/src/response"
	"app/src/service"

	"github.com/gofiber/fiber/v2"
)

type HealthCheckController struct {
	HealthCheckService service.HealthCheckService
}

func NewHealthCheckController(healthCheckService service.HealthCheckService) *HealthCheckController {
	return &HealthCheckController{
		HealthCheckService: healthCheckService,
	}
}

// @Tags Health
// @Summary Health Check
// @Description Check the status of services and database connections
// @Accept json
// @Produce json
// @Success 200 {object} example.HealthCheckResponse
// @Failure 500 {object} example.HealthCheckResponseError
// @Router /health-check [get]
func (h *HealthCheckController) Check(c *fiber.Ctx) error {
	serviceList := []response.HealthCheck{}
	isHealthy := true

	// Check database connection
	if err := h.HealthCheckService.GormCheck(); err != nil {
		isHealthy = false
		errMsg := err.Error()
		serviceList = append(serviceList, response.HealthCheck{
			Name:    constants.HealthServicePostgre,
			Status:  constants.HealthStatusDown,
			IsUp:    false,
			Message: &errMsg,
		})
	} else {
		serviceList = append(serviceList, response.HealthCheck{
			Name:   constants.HealthServicePostgre,
			Status: constants.HealthStatusUp,
			IsUp:   true,
		})
	}

	// Check memory heap
	if err := h.HealthCheckService.MemoryHeapCheck(); err != nil {
		isHealthy = false
		errMsg := err.Error()
		serviceList = append(serviceList, response.HealthCheck{
			Name:    constants.HealthServiceMemory,
			Status:  constants.HealthStatusDown,
			IsUp:    false,
			Message: &errMsg,
		})
	} else {
		serviceList = append(serviceList, response.HealthCheck{
			Name:   constants.HealthServiceMemory,
			Status: constants.HealthStatusUp,
			IsUp:   true,
		})
	}

	statusCode := fiber.StatusOK
	status := constants.HealthStatusSuccess
	if !isHealthy {
		statusCode = fiber.StatusInternalServerError
		status = constants.HealthStatusError
	}

	return c.Status(statusCode).JSON(response.HealthCheckResponse{
		Status:    status,
		Message:   constants.HealthCheckCompleted,
		Code:      statusCode,
		IsHealthy: isHealthy,
		Result:    serviceList,
	})
}
