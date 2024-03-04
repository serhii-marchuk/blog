package web

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"github.com/serhii-marchuk/blog/internal/bootstrap/configs"
	"github.com/serhii-marchuk/blog/internal/bootstrap/constructors"
	"net/http"
	"time"
)

type WebHandler struct {
	Pages       []configs.NavItem
	RedisClient *redis.Client
}

func NewWebHandler(cfg *configs.WebCfg, rc *constructors.RedisClient) *WebHandler {
	return &WebHandler{Pages: cfg.NavCfg.NavBar, RedisClient: rc.RC}
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

	return h.Render(c, http.StatusOK, page, map[string]interface{}{"pages": &h.Pages})
}

func (h *WebHandler) Render(
	c echo.Context,
	code int,
	name string,
	data map[string]interface{},
) error {
	rndr := c.Echo().Renderer
	if rndr := c.Echo().Renderer; rndr == nil {
		return errors.New("renderer not registered")
	}

	html, _ := h.RedisClient.Get(context.Background(), name).Bytes()
	if html == nil {
		buf := new(bytes.Buffer)
		rndr.Render(buf, name, data, c)
		err := h.RedisClient.Set(context.Background(), name, buf.Bytes(), time.Hour).Err()
		if err != nil {
			return err
		}

		return c.HTMLBlob(code, buf.Bytes())
	}

	return c.HTMLBlob(code, html)
}

func (h *WebHandler) Assets(c echo.Context) error {
	return c.File(fmt.Sprintf("web/assets/%s", c.Param("assetsFilePath")))
}
