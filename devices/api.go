package devices

import (
	"fmt"
	"github.com/area3001/goira/core"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"net/http"
)

func RegisterEndpoints(e *echo.Group, service *Service) {
	e.GET("", listDeviceKeys(service))
	e.GET("/:key", getDevice(service))

	e.POST("/:key/mode", setMode(service))

	e.POST("/:key/blink", blink(service))
	e.POST("/:key/fx", setFx(service))
	e.POST("/:key/rgb", sendRgb(service))

	e.DELETE("/:key", reset(service))
}

// HealthCheck godoc
// @Summary get the keys for known devices.
// @Description get the keys for known devices.
// @Tags root
// @Accept */*
// @Produce json
// @Success 200 {array} string
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
// @Success 200 {object} devices.Device
// @Param        key   path      string  true  "Device Key"
// @Router /devices/{key} [get]
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

// HealthCheck godoc
// @Summary set the mode of a device.
// @Description Set the execution mode for a device.
// @Description Execution modes define what data is allowed to be sent to/from the IRA
// @Description Valid modes are:
// @Description - 0: ExternallySet
// @Description - 1: DmxIn
// @Description - 2: DmxOut
// @Description - 3: DmxToPixelsWithIr
// @Description - 4: DmxToPixels
// @Description - 5: RgbToPixelsWithIr
// @Description - 6: RgbToPixels
// @Description - 7: FxToPixelsWithIr
// @Description - 8: FxToPixels
// @Description - 9: AutoFxWithIr
// @Description - 10: AutoFx
// @Description - 11: Emergency,
// @Tags root
// @Produce json
// @Success 200
// @Failure 400
// @Param        key   path      string  true  "Device Key"
// @Param        mode   formData      int  true  "Device Mode"
// @Router /devices/{key}/mode [post]
func setMode(service *Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		key := c.Param("key")

		mode := -1

		err := echo.FormFieldBinder(c).
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

// HealthCheck godoc
// @Summary reset a device.
// @Description Reset a device after a certain delay.
// @Tags root
// @Produce json
// @Success 200
// @Failure 400
// @Failure 500
// @Param        key   path      string  true  "Device Key"
// @Param        delay   body      int  true  "Restart delay expressed in milliseconds"
// @Router /devices/{key} [delete]
func reset(service *Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		key := c.Param("key")

		delay := 0
		err := echo.FormFieldBinder(c).
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

// HealthCheck godoc
// @Summary blink the debug led on a device.
// @Description Blink the debug led on a device for a certain number of times.
// @Tags root
// @Produce json
// @Success 200
// @Failure 400
// @Param        key   path      string  true  "Device Key"
// @Param        times   body      int  true  "The amount of times to blink"
// @Router /devices/{key}/blink [post]
func blink(service *Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		key := c.Param("key")

		times := 0
		err := echo.FormFieldBinder(c).
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

// HealthCheck godoc
// @Summary set the current effect for a device
// @Description Set the current effect for a device.
// @Description The following effects are available:
// @Description  - 0: PixelLoopFx
// @Description  - 1: RandomPixelLoopFx
// @Description  - 2: ForegroundBackgroundLoopFx
// @Description  - 3: ForegroundBackgroundSwitchFx
// @Description  - 4: Fire2021Fx
// @Tags root
// @Produce json
// @Success 200
// @Failure 400
// @Param        key   path      string  true  "Device Key"
// @Param        effect   body      core.Effect  true  "the effect to apply"
// @Router /devices/{key}/fx [post]
func setFx(service *Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		key := c.Param("key")

		var fx core.Effect
		if err := c.Bind(&fx); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error": err.Error(),
			})
		}

		if err := service.SendFx(key, &fx); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error": err.Error(),
			})
		}

		return c.NoContent(http.StatusOK)
	}
}

// HealthCheck godoc
// @Summary Send RGB Data to a device
// @Description Send RGB Data to a device
// @Tags root
// @Produce json
// @Success 200
// @Failure 400
// @Param        key   path      string  true  "Device Key"
// @Param        data   body      []byte  true  "the rgb data to send"
// @Router /devices/{key}/rgb [post]
func sendRgb(service *Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		key := c.Param("key")

		b, err := ioutil.ReadAll(c.Request().Body)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error": err.Error(),
			})
		}

		if err := service.SendRgb(key, b); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error": err.Error(),
			})
		}

		return c.NoContent(http.StatusOK)
	}
}
