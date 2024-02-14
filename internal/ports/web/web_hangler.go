package web

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

type WebHandler struct {
}

func NewWebHandler() *WebHandler {
	return &WebHandler{}
}

func (h *WebHandler) Setup(e *echo.Echo) {
	e.GET("/", h.Page)
	e.GET("/:page", h.Page)
	e.GET("/assets/:assetsFilePath", h.Assets)
}

func (h *WebHandler) Page(c echo.Context) error {
	page := c.Param("page")
	if page == "" {
		page = "home"
	}

	return c.Render(http.StatusOK, page, map[string]interface{}{"page": page})
}

func (h *WebHandler) Assets(c echo.Context) error {
	return c.File(fmt.Sprintf("web/assets/%s", c.Param("assetsFilePath")))
}
