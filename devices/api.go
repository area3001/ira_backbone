package devices

import (
	"fmt"
	"github.com/area3001/goira/core"
	"github.com/labstack/echo/v4"
	"net/http"
)

func RegisterEndpoints(e *echo.Group, service *Service) {
	e.GET("", listDeviceKeys(service))
	e.GET("/:key", getDevice(service))

	e.POST("/:key/mode", setMode(service))

	e.POST("/:key/blink", blink(service))

	e.DELETE("/:key", reset(service))
}

// HealthCheck godoc
// @Summary get the keys for known devices.
// @Description get the keys for known devices.
// @Tags root
// @Accept */*
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /devices [get]
func listDeviceKeys(service *Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		keys, err := service.Keys()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error": err.Error(),
			})
		}

		return c.JSON(http.StatusOK, keys)
	}
}

// HealthCheck godoc
// @Summary get the information for a device.
// @Description get the information for a device.
// @Tags root
// @Accept */*
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /devices/:id [get]
func getDevice(service *Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		key := c.Param("key")
		dev, err := service.GetDevice(key)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error": err.Error(),
			})
		}

		return c.JSON(http.StatusOK, dev)
	}
}

//func listDevicesHandler(service *Service) echo.HandlerFunc {
//	return func(c echo.Context) error {
//		r, err := service.ListDevices()
//		if err != nil {
//			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
//				"error": err.Error(),
//			})
//		}
//
//		return c.JSON(http.StatusOK, r)
//	}
//}

func setMode(service *Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		key := c.Param("key")

		mode := -1

		err := echo.QueryParamsBinder(c).
			Int("mode", &mode).
			BindError()

		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error": err.Error(),
			})
		}

		if mode < 0 || mode >= len(core.Modes) {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error": fmt.Sprintf("given mode %d is not between 0 and %d", mode, len(core.Modes)),
			})
		}

		if err := service.SetMode(key, core.Modes[mode]); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error": err.Error(),
			})
		}

		return c.NoContent(http.StatusOK)
	}
}

func reset(service *Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		key := c.Param("key")

		delay := 0
		err := echo.QueryParamsBinder(c).
			Int("delay", &delay).
			BindError()

		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error": err.Error(),
			})
		}

		if delay < 0 {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error": "the delay should be larger or equal to 0",
			})
		}

		if err := service.Reset(key, delay); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error": err.Error(),
			})
		}

		return c.NoContent(http.StatusOK)
	}
}

func blink(service *Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		key := c.Param("key")

		times := 0
		err := echo.QueryParamsBinder(c).
			Int("times", &times).
			BindError()

		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error": err.Error(),
			})
		}

		if times < 1 {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error": "the amount of times to flash should be larger or equal to 1",
			})
		}

		if err := service.Blink(key, times); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error": err.Error(),
			})
		}

		return c.NoContent(http.StatusOK)
	}
}
