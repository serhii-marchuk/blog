package web

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/serhii-marchuk/blog/internal/bootstrap/configs"
	"net/http"
)

type WebHandler struct {
	Pages []configs.NavItem
}

func NewWebHandler(cfg *configs.WebCfg) *WebHandler {
	return &WebHandler{Pages: cfg.NavCfg.NavBar}
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

	for k, item := range h.Pages {
		if item.Name == page {
			h.Pages[k].Active = true
		} else {
			h.Pages[k].Active = false
		}
	}

	return c.Render(http.StatusOK, page, map[string]interface{}{"pages": &h.Pages})
}

func (h *WebHandler) Assets(c echo.Context) error {
	return c.File(fmt.Sprintf("web/assets/%s", c.Param("assetsFilePath")))
}
