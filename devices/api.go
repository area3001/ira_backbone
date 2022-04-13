package devices

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func RegisterEndpoints(e *echo.Group, service *Service) {
	e.GET("", func(c echo.Context) error {
		r, err := service.ListDevices()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error": err.Error(),
			})
		}

		return c.JSON(http.StatusOK, r)
	})
}
