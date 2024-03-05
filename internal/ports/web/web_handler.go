package web

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/docker/distribution/uuid"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"github.com/serhii-marchuk/blog/internal/bootstrap/configs"
	"github.com/serhii-marchuk/blog/internal/bootstrap/constructors"
	"github.com/serhii-marchuk/blog/internal/models"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type WebHandler struct {
	Pages       []configs.NavItem
	RedisClient *redis.Client
	Db          *gorm.DB
}

func NewWebHandler(
	cfg *configs.WebCfg,
	rc *constructors.RedisClient,
	db *constructors.Db,
) *WebHandler {
	return &WebHandler{Pages: cfg.NavCfg.NavBar, RedisClient: rc.RC, Db: db.DB}
}

func (h *WebHandler) Setup(e *echo.Echo) {
	e.GET("/", h.Page)
	e.GET("/:page", h.Page)
	e.POST("/contact/form", h.HandleContactForm)
	e.GET("/assets/:assetsFilePath", h.Assets)
}

func (h *WebHandler) Page(c echo.Context) error {
	page := c.Param("page")
	if page == "" {
		page = "home"
	}

	return h.Render(c, http.StatusOK, page)
}

func (h *WebHandler) Render(
	c echo.Context,
	code int,
	name string,
) error {
	rndr := c.Echo().Renderer
	if rndr := c.Echo().Renderer; rndr == nil {
		return errors.New("renderer not registered")
	}

	html, _ := h.RedisClient.Get(context.Background(), name).Bytes()
	if html == nil {

		for k, item := range h.Pages {
			if item.Name == name {
				h.Pages[k].Active = true
			} else {
				h.Pages[k].Active = false
			}
		}

		buf := new(bytes.Buffer)
		rndr.Render(buf, name, map[string]interface{}{"pages": &h.Pages}, c)
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

func (h *WebHandler) HandleContactForm(c echo.Context) error {
	h.Db.Save(&models.ContactRequest{
		ID:          uuid.Generate(),
		FirstName:   c.FormValue("first_name"),
		LastName:    c.FormValue("last_name"),
		Email:       c.FormValue("email"),
		Description: c.FormValue("description"),
		Status:      "new",
	})

	return h.Render(c, http.StatusOK, "contact")
}
